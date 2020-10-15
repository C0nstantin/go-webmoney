package webmoney

import "encoding/xml"

type X21Request struct {
	merchantRequest
	LmiPayeePurse       string `xml:"lmi_payee_purse"`
	LmiDayLimit         string `xml:"lmi_day_limit"`
	LmiWeekLimit        string `xml:"lmi_week_limit"`
	LmiMonthLimit       string `xml:"lmi_month_limit"`
	LmiClientnumber     string `xml:"lmi_clientnumber"`
	LmiClientnubmerType string `xml:"lmi_clientnumber_type"`
	LmiSmsType          string `xml:"lmi_sms_type"`
	Lang                string `xml:"lang"`
}

func (x X21Request) GetSignSource(s string) (string, error) {
	return x.LmiPayeePurse + x.LmiClientnumber + x.LmiClientnubmerType + x.LmiSmsType, nil

}

type X21Response struct {
	merchantResponse
	Slavepurse  string   `xml:"slavepurse"`
	Slavewmid   string   `xml:"slavewmid"`
	Smssecureid string   `xml:"smssecureid"`
	TrustX21    TrustX21 `xml:"trust"`
}
type TrustX21 struct {
	XMLName     xml.Name `xml:"trust"`
	PurseId     string   `xml:"purseid,attr"`
	RealsmsType string   `xml:"realsmstype"`
}

type X212Request struct {
	merchantRequest
	LmiPurseId          string `xml:"lmi_purseid"`
	LmiClientnumberCode string `xml:"lmi_clientnumber_code"`
}
type TrustX212 struct {
	XMLName    xml.Name `xml:"trust"`
	Id         string   `xml:"id,attr"`
	Slavepurse string   `xml:"slavepurse"`
	Slavewmid  string   `xml:"slavewmid"`
	Masterwmid string   `xml:"masterwmid"`
}
type X212Response struct {
	merchantResponse
	TrustX212    TrustX212 `xml:"trust"`
	Smssentstate string    `xml:"smssentstate"`
}

func (x X212Request) GetSignSource(s string) (string, error) {
	return x.LmiPurseId + x.LmiClientnumberCode, nil
}

func (w *WmClient) TrustRequest(x X21Request) (X21Response, error) {

	X := Merchant{
		Request:   &x,
		Interface: XInterface{Name: "TrustRequest", Type: "merchant"},
		Client:    w,
	}

	result := X21Response{}
	err := X.getResult(&result)

	return result, err
}
func (w *WmClient) TrustConfirm(x X212Request) (X212Response, error) {

	X := Merchant{
		Request:   &x,
		Interface: XInterface{Name: "TrustConfirm", Type: "merchant"},
		Client:    w,
	}

	result := X212Response{}
	err := X.getResult(&result)

	return result, err
}
