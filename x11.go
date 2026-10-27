package webmoney

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

var (
	ErrAttestatRecalled = errors.New("attestat recalled")
	ErrAttestatNotFound = errors.New("attestat not found")

	// ErrPassportUnavailable means the passport service did not give a
	// trustworthy answer: transport error or timeout, non-200 HTTP status,
	// unparsable body, missing or non-numeric retval attribute, or a
	// retval="0" response without <certinfo> for the requested WMID.
	// Callers must not treat it as "WMID not found".
	ErrPassportUnavailable = errors.New("passport unavailable")

	// ErrPassportRetval means the passport answered with retval != 0 other
	// than retvalInvalidWmid. The error text carries retval and retdesc.
	ErrPassportRetval = errors.New("passport returned error retval")
)

// retvalInvalidWmid is the X11 answer to a malformed passportwmid
// ("неверно указан проверяемый WMID"). It is a definite answer, not an
// outage: such a WMID does not exist.
const retvalInvalidWmid = 2

// passportError is ErrPassportUnavailable with the underlying cause attached,
// so errors.Is(err, ErrPassportUnavailable) and errors.As(err, &netErr) both work.
type passportError struct {
	msg   string
	cause error
}

func (e *passportError) Error() string {
	if e.cause != nil {
		return ErrPassportUnavailable.Error() + ": " + e.msg + ": " + e.cause.Error()
	}
	return ErrPassportUnavailable.Error() + ": " + e.msg
}

func (e *passportError) Unwrap() error { return e.cause }

func (e *passportError) Is(target error) bool { return target == ErrPassportUnavailable }

func unavailable(msg string, cause error) error {
	return &passportError{msg: msg, cause: cause}
}

// WmidStatus is the outcome of CheckWmid.
type WmidStatus int

const (
	// WmidUnknown: the passport gave no trustworthy answer (see the error).
	WmidUnknown WmidStatus = iota
	// WmidExists: the passport knows this WMID (it returned <attestat> or <wmids>).
	WmidExists
	// WmidNotFound: the passport answered retval="0" but knows nothing about
	// the WMID — it was deleted from WebMoney or never existed.
	WmidNotFound
)

func (s WmidStatus) String() string {
	switch s {
	case WmidExists:
		return "exists"
	case WmidNotFound:
		return "not found"
	default:
		return "unknown"
	}
}

type RequestX11 struct {
	XMLName      xml.Name `xml:"request"`
	Wmid         string   `xml:"wmid"`
	PassportWmid string   `xml:"passportwmid"`
	Sign         string   `xml:"sign"`
	Dict         string   `xml:"params>dict"`
	Info         string   `xml:"params>info"`
	Mode         string   `xml:"params>mode"`
}

type ResponseX11 struct {
	XMLName xml.Name `xml:"response"`
	// Retval is guaranteed to come from an explicit retval attribute when the
	// struct was filled by this package: responses without it are rejected
	// with ErrPassportUnavailable before decoding.
	Retval   int      `xml:"retval,attr"`
	Retdesc  string   `xml:"retdesc,attr"`
	CertInfo CertInfo `xml:"certinfo"`
}

type CertInfo struct {
	XMLName  xml.Name `xml:"certinfo"`
	Attestat Attestat `xml:"attestat>row"`
	Wmid     string   `xml:"wmid,attr"`
	Wmids    []WmidP  `xml:"wmids>row"`
	UserInfo UserInfo `xml:"userinfo>value>row"`
}

type WmidP struct {
	XMLName     xml.Name `xml:"row"`
	Wmid        string   `xml:"wmid,attr"`
	Info        string   `xml:"info,attr"`
	Nickname    string   `xml:"nickname,attr"`
	Datereg     string   `xml:"datereg,attr"`
	Ctype       string   `xml:"ctype,attr"`
	Companyname string   `xml:"companyname,attr"`
	Companyid   string   `xml:"companyid,attr"`
}

type UserInfo struct {
	XMLName  xml.Name `xml:"row"`
	Nickname string   `xml:"nickname,attr"`
	Country  string   `xml:"country,attr"`
	City     string   `xml:"city,attr"`
	Zipcode  string   `xml:"zipcode,attr"`
	Adres    string   `xml:"adres,attr"`
	Fname    string   `xml:"fname,attr"`
	Iname    string   `xml:"iname,attr"`
	Oname    string   `xml:"oname,attr"`
	Pnomer   string   `xml:"pnomer,attr"`
	Pdate    string   `xml:"pdate,attr"`
	Pdateend string   `xml:"pdateend,attr"`
	Pcountry string   `xml:"pcountry,attr"`
	Pbywhom  string   `xml:"pbywhom,attr"`
	Inn      string   `xml:"inn,attr"`
	Bplace   string   `xml:"bplace,attr"`
	Bday     string   `xml:"bday,attr"`
	Bmonth   string   `xml:"bmonth,attr"`
	Byear    string   `xml:"byear,attr"`
	Icq      string   `xml:"icq,attr"`
	Email    string   `xml:"email,attr"`
	Web      string   `xml:"web,attr"`
	Phone    string   `xml:"phone,attr"`
	CapOwner string   `xml:"cap_owner,attr"`
	Pasdoc   string   `xml:"pasdoc,attr"`
	Regdoc   string   `xml:"regdoc,attr"`
	Inndoc   string   `xml:"inndoc,attr"`
	Photoid  string   `xml:"phoneid,attr"`
	Jabberid string   `xml:"jabberid,attr"`
	Sex      string   `xml:"sex,attr"`
	Infoopen string   `xml:"infoopen,attr"`
}

type Attestat struct {
	TID        string `xml:"tid,attr"`
	Recalled   string `xml:"recalled,attr"`
	Datecrt    string `xml:"datecrt,attr"`
	Dateupd    string `xml:"dateupd,attr"`
	Regnikname string `xml:"regnickname,attr"`
	Regwmid    string `xml:"regwmid,attr"`
	Status     string `xml:"status,attr"`
	Notary     string `xml:"notary,attr"`
}

func getPassportDomain() string {
	apiPassportDomain := "apipassport.web.money"
	if v, ok := os.LookupEnv("API_PASSPORT_DOMAIN"); ok && isValidDomain(v) {
		apiPassportDomain = v
	}
	return apiPassportDomain
}

func getPassportURL() string {
	return fmt.Sprintf("https://%s/asp/XMLGetWMPassport.asp", getPassportDomain())
}

// maxPassportResponseSize limits the passport response body. A real X11
// answer is a few kilobytes; 1 MB leaves a wide margin.
const maxPassportResponseSize = 1 << 20

var (
	passportMu      sync.Mutex
	passportClient  *http.Client
	passportCustom  bool
	passportTimeout time.Duration
)

// passportTimeoutFromEnv reads API_PASSPORT_TIMEOUT (a Go duration such as
// "15s"). Invalid or non-positive values fall back to defaultHTTPTimeout.
func passportTimeoutFromEnv() time.Duration {
	if v, ok := os.LookupEnv("API_PASSPORT_TIMEOUT"); ok {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return defaultHTTPTimeout
}

// SetPassportTimeout sets the HTTP timeout used by the package-level X11
// functions (GetInfoWmid, WMPassportByWmid, IssetWmid, CheckWmid).
// d <= 0 restores the default (API_PASSPORT_TIMEOUT or 30s).
// It does not touch a client installed with SetPassportHTTPClient.
func SetPassportTimeout(d time.Duration) {
	passportMu.Lock()
	defer passportMu.Unlock()
	if d <= 0 {
		d = 0
	}
	passportTimeout = d
	if !passportCustom {
		passportClient = nil
	}
}

// SetPassportHTTPClient replaces the HTTP client used by the package-level
// X11 functions, e.g. to add a proxy or a custom transport. The caller is
// responsible for the client's Timeout. nil restores the default client.
func SetPassportHTTPClient(c *http.Client) {
	passportMu.Lock()
	defer passportMu.Unlock()
	passportClient = c
	passportCustom = c != nil
}

// getPassportClient returns the cached package-level client with a timeout.
func getPassportClient() *http.Client {
	passportMu.Lock()
	defer passportMu.Unlock()
	if passportClient == nil {
		t := passportTimeout
		if t <= 0 {
			t = passportTimeoutFromEnv()
		}
		passportClient = &http.Client{
			Timeout: t,
			// An API endpoint must not redirect; following one (as GET, to
			// another host) would only produce a response that is not from
			// the passport. Return the 3xx as is: the status check rejects it.
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
		passportCustom = false
	}
	return passportClient
}

// passportEnvelope is what a response must carry to be trusted: an explicit
// retval attribute and a <certinfo> about the WMID that was asked for.
type passportEnvelope struct {
	XMLName xml.Name `xml:"response"`
	// *string, not int: encoding/xml decodes a missing or empty attribute
	// into int 0, which would look like an explicit success.
	Retval   *string `xml:"retval,attr"`
	CertInfo *struct {
		Wmid string `xml:"wmid,attr"`
	} `xml:"certinfo"`
}

func decodePassportXML(body []byte, target interface{}) error {
	dec := xml.NewDecoder(bytes.NewReader(body))
	dec.CharsetReader = charset.NewReader
	return dec.Decode(target)
}

// postPassportRequest sends an X11 request for wmid and decodes the response
// into target only if the response is trustworthy: HTTP 200, well-formed XML
// with a numeric retval attribute and, for retval="0", a <certinfo> whose
// wmid attribute echoes the requested WMID. Anything else is
// ErrPassportUnavailable. A response with retval != 0 is decoded and returned
// without error — the caller decides what it means.
func postPassportRequest(client *http.Client, wmid string, target *ResponseX11) error {
	out, err := xml.MarshalIndent(RequestX11{PassportWmid: wmid}, " ", "  ")
	if err != nil {
		return err
	}

	resp, err := client.Post(getPassportURL(), "text/xml", bytes.NewReader(out))
	if err != nil {
		return unavailable("request failed", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxPassportResponseSize))
	if err != nil {
		return unavailable("reading body", err)
	}
	if resp.StatusCode != http.StatusOK {
		return unavailable(fmt.Sprintf("HTTP %d", resp.StatusCode), nil)
	}

	var env passportEnvelope
	if err := decodePassportXML(body, &env); err != nil {
		return unavailable("malformed response", err)
	}
	if env.Retval == nil {
		return unavailable("response has no retval attribute", nil)
	}
	retval, err := strconv.Atoi(strings.TrimSpace(*env.Retval))
	if err != nil {
		return unavailable(fmt.Sprintf("retval %q is not a number", truncate(*env.Retval)), nil)
	}

	if err := decodePassportXML(body, target); err != nil {
		return unavailable("malformed response", err)
	}
	target.Retval = retval

	if retval != 0 {
		return nil
	}
	// certinfo.wmid echoes the request. It does not prove the WMID exists,
	// but it does prove the answer is about the WMID we asked for.
	if env.CertInfo == nil || env.CertInfo.Wmid != wmid {
		return unavailable("response is not about wmid "+wmid, nil)
	}
	return nil
}

// maxRetdescLen bounds retdesc in error texts: they end up in logs.
const maxRetdescLen = 200

func truncate(s string) string {
	if len(s) <= maxRetdescLen {
		return s
	}
	return s[:maxRetdescLen] + "..."
}

// retvalError builds ErrPassportRetval for a response with retval != 0.
func retvalError(v ResponseX11) error {
	if v.Retdesc != "" {
		return fmt.Errorf("%w: retval=%d %s", ErrPassportRetval, v.Retval, truncate(v.Retdesc))
	}
	return fmt.Errorf("%w: retval=%d", ErrPassportRetval, v.Retval)
}

// classifyWmid decides whether the passport knows the WMID. Only an
// <attestat> or a <wmids> entry proves that; certinfo.wmid is an echo of the
// request and a deleted or made-up WMID yields the same empty <certinfo>.
func classifyWmid(v ResponseX11) WmidStatus {
	if v.Retval != 0 {
		return WmidUnknown
	}
	known := v.CertInfo.Attestat.TID != "" ||
		v.CertInfo.Attestat.Datecrt != "" ||
		len(v.CertInfo.Wmids) > 0
	if known {
		return WmidExists
	}
	return WmidNotFound
}

func checkWmid(client *http.Client, wmid string) (WmidStatus, error) {
	var v ResponseX11
	if err := postPassportRequest(client, wmid, &v); err != nil {
		return WmidUnknown, err
	}
	switch v.Retval {
	case 0:
		return classifyWmid(v), nil
	case retvalInvalidWmid:
		return WmidNotFound, nil
	default:
		return WmidUnknown, retvalError(v)
	}
}

func getInfoWmid(client *http.Client, wmid string) (ResponseX11, error) {
	var v ResponseX11
	if err := postPassportRequest(client, wmid, &v); err != nil {
		return ResponseX11{}, err
	}
	switch v.Retval {
	case 0:
		return checkAttestat(v)
	case retvalInvalidWmid:
		// A malformed WMID has no passport. Kept as ErrAttestatNotFound for
		// backward compatibility: callers match this text.
		return v, ErrAttestatNotFound
	default:
		return v, retvalError(v)
	}
}

func wmPassportByWmid(client *http.Client, wmid string) (*ResponseX11, error) {
	v := &ResponseX11{}
	if err := postPassportRequest(client, wmid, v); err != nil {
		return nil, err
	}
	return v, nil
}

func checkAttestat(v1 ResponseX11) (ResponseX11, error) {
	if v1.CertInfo.Attestat.Recalled == "1" {
		return v1, ErrAttestatRecalled
	}
	if v1.CertInfo.Attestat.TID != "" {
		return v1, nil
	}
	return v1, ErrAttestatNotFound
}

// CheckWmid asks the passport whether a WMID exists and returns one of three
// outcomes:
//   - WmidExists, nil     — the passport knows the WMID (attestat or wmids present);
//   - WmidNotFound, nil   — retval="0" but the passport knows nothing about it
//     (deleted from WebMoney or never existed), or retval="2" (malformed WMID);
//   - WmidUnknown, err    — no trustworthy answer: ErrPassportUnavailable
//     (timeout, HTTP error, broken XML, no retval) or ErrPassportRetval
//     (retval != 0). Do not act on it.
//
// Uses a package-level client with a timeout (see SetPassportTimeout).
func CheckWmid(wmid string) (WmidStatus, error) {
	return checkWmid(getPassportClient(), wmid)
}

// IssetWmid reports whether a WMID exists. It is CheckWmid folded into a bool:
// (true, nil) for WmidExists, (false, nil) for WmidNotFound and (false, err)
// when the passport gave no trustworthy answer. Callers that ignore the error
// will treat a passport outage as "WMID does not exist" — prefer CheckWmid.
func IssetWmid(wmid string) (bool, error) {
	st, err := CheckWmid(wmid)
	return st == WmidExists, err
}

// GetInfoWmid returns the passport info for a WMID and checks its attestat:
// ErrAttestatRecalled if the attestat is recalled, ErrAttestatNotFound if the
// response carries no attestat (this includes WMIDs the passport does not
// know and malformed WMIDs, retval="2" — use CheckWmid to tell them apart).
// ErrPassportUnavailable and ErrPassportRetval mean the passport gave no
// trustworthy answer; the returned ResponseX11 then carries Retval/Retdesc.
//
// Uses a package-level client with a timeout (see SetPassportTimeout).
// Prefer WmClient.GetInfoWmid() when a WmClient is at hand.
func GetInfoWmid(wmid string) (ResponseX11, error) {
	return getInfoWmid(getPassportClient(), wmid)
}

// WMPassportByWmid returns the raw passport response. The transport is strict
// (HTTP 200, numeric retval, certinfo for the requested WMID — otherwise
// ErrPassportUnavailable), but a response with retval != 0 is returned as is
// so the caller can inspect Retval and Retdesc.
//
// Uses a package-level client with a timeout (see SetPassportTimeout).
func WMPassportByWmid(wmid string) (*ResponseX11, error) {
	return wmPassportByWmid(getPassportClient(), wmid)
}

// CheckWmid is CheckWmid using the WmClient's cached HTTP client.
func (w *WmClient) CheckWmid(wmid string) (WmidStatus, error) {
	client, err := w.getHTTPClient()
	if err != nil {
		return WmidUnknown, err
	}
	return checkWmid(client, wmid)
}

// IssetWmid is IssetWmid using the WmClient's cached HTTP client.
func (w *WmClient) IssetWmid(wmid string) (bool, error) {
	st, err := w.CheckWmid(wmid)
	return st == WmidExists, err
}

// GetInfoWmid is GetInfoWmid using the WmClient's cached HTTP client with
// proper TLS and timeout settings.
func (w *WmClient) GetInfoWmid(wmid string) (ResponseX11, error) {
	client, err := w.getHTTPClient()
	if err != nil {
		return ResponseX11{}, err
	}
	return getInfoWmid(client, wmid)
}

// WMPassportByWmid is WMPassportByWmid using the WmClient's cached HTTP client.
func (w *WmClient) WMPassportByWmid(wmid string) (*ResponseX11, error) {
	client, err := w.getHTTPClient()
	if err != nil {
		return nil, err
	}
	return wmPassportByWmid(client, wmid)
}
