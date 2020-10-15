package webmoney

import "encoding/xml"

type Paymenttags struct {
	LmiPayeePurse        string `xml:"lmi_payee_purse"`
	LmiPaymentNo         string `xml:"lmi_payment_no"`
	LmiPaymentNoType     string `xml:"lmi_payment_no_type"`
	LmiPaymentAmount     string `xml:"lmi_payment_amount"`
	LmiPaymentDesc       string `xml:"lmi_payment_desc"`
	LmiPaymentDescBase64 string `xml:"lmi_payment_desc_base64"`
	LmiClientnumber      string `xml:"lmi_clientnumber"`
	LmiClientnubmerType  string `xml:"lmi_clientnumber_type"`
	LmiSmsType           string `xml:"lmi_sms_type"`
	LmiShopId            string `xml:"lmi_shop_id"`
	LmiHold              string `xml:"LMI_HOLD"`
	Lang                 string `xml:"lang"`
	EnulatedFlag         string `xml:"emulated_flag"`
}

type X22Request struct {
	XMLName               xml.Name    `xml:"merchant.request"`
	Wmid                  string      `xml:"signtags>wmid"`
	Sign                  string      `xml:"signtags>sign"`
	Sha256                string      `xml:"signtags>sha256"`
	Md5                   string      `xml:"signtags>md5"`
	SecretKey             string      `xml:"signtags>secret_key"`
	Lang                  string      `xml:"signtags>lang"`
	Validityperiodinhours string      `xml:"signtags>validityperiodinhours"`
	Paymenttags           Paymenttags `xml:"paymenttags"`
}

func (x *X22Request) setWmid(s string) {
	x.Wmid = s

}

func (x *X22Request) setSign(s string) {
	x.Sign = s

}

func (x *X22Request) setSha256(s string) {
	x.Sha256 = s

}

func (x *X22Request) GetSignSource(s string) (string, error) {
	return x.Paymenttags.LmiPayeePurse + x.Paymenttags.LmiPaymentNo + x.Validityperiodinhours, nil
}

type X22Response struct {
	merchantResponse
	Transtoken            string `xml:"transtoken"`
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
