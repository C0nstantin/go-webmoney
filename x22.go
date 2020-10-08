package webmoney

import "encoding/xml"

type X22Request struct {
	XMLName               xml.Name   `xml:"merchant.request"`
	Wmid                  string     `xml:"signtags>wmid"`
	Sign                  string     `xml:"signtags>sign"`
	Sha256                string     `xml:"signtags>sha256"`
	Md5                   string     `xml:"signtags>md5"`
	SecretKey             string     `xml:"signtags>secret_key"`
	Lang                  string     `xml:"signtags>lang"`
	Validityperiodinhours string     `xml:"signtags>validityperiodinhours"`
	Paymenttags           X20Request `xml:"paymenttags"`
}

func (x X22Request) setWmid(s string) {
	x.Wmid = s
	//panic("implement me")
}

func (x X22Request) setSign(s string) {
	x.Sign = s
	//panic("implement me")
}

func (x X22Request) setSha256(s string) {
	x.Sha256 = s
	//panic("implement me")
}

func (x X22Request) GetSignSource(s string) (string, error) {
	return x.Wmid + x.Paymenttags.LmiPayeePurse + x.Paymenttags.LmiPaymentNo + x.Validityperiodinhours, nil
}

type X22Response struct {
	TransToken            string `xml:"transtoken"`
	Validityperiodinhours string `xml:"validityperiodinhours"`
}

func (w *WmClient) TransSave(x X22Request) (X22Response, error) {
	X := Merchant{
		Request:   &x,
		Interface: XInterface{Name: "TransSave", Type: "merchant"},
		Client:    w,
	}
	result := X22Response{}
	err := X.getResult(&result)
	return result, err
}
