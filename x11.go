package webmoney

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	ErrAttestatRecalled = errors.New("attestat recalled")
	ErrAttestatNotFound = errors.New("attestat not found")
)

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
	XMLName  xml.Name `xml:"response"`
	Retval   int      `xml:"retval,attr"`
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

// postPassportRequest sends the XML request body using the given http.Client
// and decodes the response into the provided target.
func postPassportRequest(client *http.Client, xmlBody string, target interface{}) error {
	resp, err := client.Post(getPassportURL(), "text/xml", strings.NewReader(xmlBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return err
	}

	r := bytes.NewReader(body)
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReader
	return dec.Decode(target)
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

func IssetWmid(wmid string) (bool, error) {
	result, err := GetInfoWmid(wmid)
	if result.CertInfo.Wmid == "" || err != nil {
		return false, err
	}
	return true, nil
}

// GetInfoWmid is a standalone (package-level) function that uses http.DefaultClient.
// Prefer WmClient.GetInfoWmid() for proper transport and timeout handling.
func GetInfoWmid(wmid string) (ResponseX11, error) {
	v := RequestX11{PassportWmid: wmid}
	v1 := ResponseX11{}
	out, err := xml.MarshalIndent(v, " ", "  ")
	if err != nil {
		return v1, err
	}

	if err := postPassportRequest(http.DefaultClient, string(out), &v1); err != nil {
		return ResponseX11{}, err
	}

	return checkAttestat(v1)
}

// WMPassportByWmid is a standalone (package-level) function that uses http.DefaultClient.
// Prefer WmClient.WMPassportByWmid() for proper transport and timeout handling.
func WMPassportByWmid(wmid string) (*ResponseX11, error) {
	v := RequestX11{PassportWmid: wmid}

	out, err := xml.MarshalIndent(v, " ", "  ")
	if err != nil {
		return nil, err
	}

	v1 := &ResponseX11{}
	if err := postPassportRequest(http.DefaultClient, string(out), v1); err != nil {
		return nil, err
	}
	return v1, nil
}

// GetInfoWmid uses the WmClient's cached HTTP client with proper TLS and timeout settings.
func (w *WmClient) GetInfoWmid(wmid string) (ResponseX11, error) {
	client, err := w.getHTTPClient()
	if err != nil {
		return ResponseX11{}, err
	}

	v := RequestX11{PassportWmid: wmid}
	v1 := ResponseX11{}
	out, err := xml.MarshalIndent(v, " ", "  ")
	if err != nil {
		return v1, err
	}

	if err := postPassportRequest(client, string(out), &v1); err != nil {
		return ResponseX11{}, err
	}

	return checkAttestat(v1)
}

func (w *WmClient) IssetWmid(wmid string) (bool, error) {
	result, err := w.GetInfoWmid(wmid)
	if result.CertInfo.Wmid == "" || err != nil {
		return false, err
	}
	return true, nil
}
