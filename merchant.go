package webmoney

import (
	"bytes"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"

	"github.com/C0nstantin/go-webmoney/wmsigner"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

type SignedM interface {
	Signed
	setWmid(string)
	setSign(string)
	setSha256(string)
}
type RetReturn interface {
	GetRetDesc() string
	GetRetVal() int64
}

type Merchant struct {
	Request   SignedM //X18Request//merchantRequest
	Result    RetReturn
	ResultStr string
	Interface XInterface
	Client    *WmClient
}

func (m *Merchant) getResult(result RetReturn) error {
	str, err := m.sendRequest()
	if err != nil {
		return err
	}
	if err := m.parseResponse(result, str); err != nil {
		return err
	}
	if m.Result.GetRetVal() != 0 {
		err := errors.New(strconv.FormatInt(m.Result.GetRetVal(), 10) + "   " + m.Result.GetRetDesc())
		return err

	}
	return nil
}

func (m *Merchant) parseResponse(resp RetReturn, responseStr string) error {
	v := resp
	DebugLog("merchant response: ", responseStr)
	r := bytes.NewReader([]byte(responseStr))
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReader
	err := dec.Decode(&v)
	m.Result = v
	if err != nil {
		return err
	}
	return nil
}

func (m *Merchant) sendRequest() (string, error) {

	url := m.Interface.GetUrl(m.Client.IsClassic())

	v := m.Request
	v.setWmid(m.Client.Wmid)

	if m.Client.IsClassic() {

		signer := wmsigner.NewSigner(m.Client.Wmid, m.Client.Pass, m.Client.Key)
		str, err := m.Request.GetSignSource("")
		if err != nil {
			return "", err
		}
		str = m.Client.Wmid + str
		sign, err := signer.Sign(str)
		if err != nil {
			return "", nil
		}
		v.setSign(sign)
	} else {
		str, err := m.Request.GetSignSource(m.Client.SecretKey)
		if err != nil {
			return "", err
		}
		str = m.Client.Wmid + str
		sign := fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
		v.setSha256(sign)
	}
	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		return "", err
	}
	body := "<?xml version=\"1.0\" encoding=\"utf-8\"?> \n" + string(output)
	DebugLog(" merchant request: to " + url + "\n" + body)
	return m.Client.sendRequest(url, body)
}
