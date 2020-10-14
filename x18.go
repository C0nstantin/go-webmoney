package webmoney

import "encoding/xml"

type X18Request struct {
	merchantRequest
	LmiPayeePurse    string `xml:"lmi_payee_purse"`
	LmiPaymentNo     string `xml:"lmi_payment_no"`
	LmiPaymentNoType string `xml:"lmi_payment_no_type"`
}

func (x X18Request) GetSignSource(s string) (string, error) {
	return x.LmiPayeePurse + x.LmiPaymentNo + s, nil
}

type MerchantOperation struct {
	XmlName        xml.Name `xml:"operation"`
	Wmtransid      string   `xml:"wmtransid,attr"`
	Wminvoiceid    string   `xml:"wminvoiceid,attr"`
	Realsmstype    string   `xml:"realsmstype"`
	Amount         string   `xml:"amount"`
	Operdate       string   `xml:"operdate"`
	Purpose        string   `xml:"purpose"`
	Purposefrom    string   `xml:"pursefrom"`
	Wmidfrom       string   `xml:"wmidfrom"`
	HoldPeriod     string   `xml:"hold_period"`
	HoldState      string   `xml:"hold_state"`
	Capitallerflag string   `xml:"capitallerflag"`
	Enumflag       string   `xml:"enumflag"`
	IPAddress      string   `xml:"IPAddress"`
	TelepayPhone   string   `xml:"telepat_phone"`
	TelepayPaytype string   `xml:"telepat_paytype"`
	PaymerNumber   string   `xml:"paymer_number"`
	PaymerEmail    string   `xml:"paymer_email"`
	PaymerType     string   `xml:"paymer_type"`
	CashierNumber  string   `xml:"cashier_number"`
	CashierDate    string   `xml:"cashier_date"`
	CashierAmount  string   `xml:"cashier_amount"`
	SdpType        string   `xml:"sdp_type"`
}

func (w *WmClient) TransGet(m X18Request) (MerchantOperation, error) {

	X := Merchant{
		Request:   &m,
		Interface: XInterface{Name: "TransGet", Type: "merchant"},
		Client:    w,
	}
	result := struct {
		merchantResponse
		Operation MerchantOperation
	}{}
	err := X.getResult(&result)

	return result.Operation, err
}
