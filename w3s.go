package webmoney

import (
	"bytes"
	"encoding/xml"
	"errors"
	"strconv"

	"github.com/C0nstantin/go-webmoney/wmsigner"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

type Signed interface {
	GetSignSource(string) (string, error)
}

type W3s struct {
	Request   Signed
	Result    responseW3s
	ResultStr string
	Interface XInterface
	Client    *WmClient
}

func (w3s *W3s) getResult(result interface{}) error {
	str, err := w3s.sendRequest()
	if err != nil {
		return err
	}
	if err := w3s.parseResponse(result, str); err != nil {
		return err
	}
	if w3s.Result.Retval != 0 {
		err := errors.New(strconv.FormatInt(w3s.Result.Retval, 10) + "   " + w3s.Result.Retdesc)
		return err
	}
	return nil
}

func (w3s *W3s) parseResponse(resp interface{}, responseStr string) error {
	v := responseW3s{
		Response: resp,
	}
	r := bytes.NewReader([]byte(responseStr))
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReader
	err := dec.Decode(&v)
	w3s.Result = v
	if err != nil {
		return err
	}
	return nil
}

func (w3s *W3s) sendRequest() (string, error) {
	reqn := Reqn()
	v := &requestW3s{
		Wmid: w3s.Client.Wmid,
		Reqn: reqn,
	}

	if w3s.Client.IsClassic() {

		s := wmsigner.NewSigner(w3s.Client.Wmid, w3s.Client.Pass, w3s.Client.Key)
		str, err := w3s.Request.GetSignSource(reqn)

		if w3s.Interface.Name == "ClassicAuth" || w3s.Interface.Name == "TrustSave2" {
			str = w3s.Client.Wmid + str
		}

		if err != nil {
			return "", err
		}

		if result, err := s.Sign(str); err != nil {
			return "", err
		} else {
			v.Sign = result
		}

	} else {

		v.Sign = ""

	}
	url := w3s.Interface.GetUrl(w3s.Client.IsClassic())
	v.Request = w3s.Request
	output, err := xml.MarshalIndent(v, "  ", "    ")
	if err != nil {
		return "", err
	}
	body := "<?xml version=\"1.0\" encoding=\"utf-8\"?> \n" + string(output)
	return w3s.Client.sendRequest(url, body)
}
