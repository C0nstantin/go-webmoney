package webmoney

type X20Request struct {
	merchantRequest
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

func (x X20Request) GetSignSource(s string) (string, error) {
	return x.LmiPayeePurse + x.LmiPaymentNo + x.LmiClientnumber + x.LmiClientnubmerType, nil
}

type X202Request struct {
	merchantRequest
	LmiPayeePurse       string `xml:"lmi_payee_purse"`
	LmiClientnumberCode string `xml:"lmi_clientnumber_code"`
	LmiWminvoiceid      string `xml:"lmi_wminvoiceid"`
}

func (x X202Request) GetSignSource(s string) (string, error) {
	return x.LmiPayeePurse + x.LmiWminvoiceid + x.LmiClientnumberCode, nil
}

func (w *WmClient) TransRequest(x X20Request) (MerchantOperation, error) {
	X := Merchant{
		Request:   &x,
		Interface: XInterface{Name: "TransRequest", Type: "merchant"},
		Client:    w,
	}
	result := MerchantOperation{}
	err := X.getResult(&result)

	return result, err
}

func (w *WmClient) TransConfirm(x X202Request) (MerchantOperation, error) {

	X := Merchant{
		Request:   &x,
		Interface: XInterface{Name: "TransConfirm", Type: "merchant"},
		Client:    w,
	}

	result := MerchantOperation{}
	err := X.getResult(&result)

	return result, err
}
