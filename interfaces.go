package webmoney

import "encoding/xml"

type XInterface struct {
	Name string
	Type string
}

func (x *XInterface) GetUrl(isClassic bool) string {
	var url string
	if x.Type == "w3s" || x.Type == "passport" {
		if isClassic {
			url = "https://w10s.webmoney.ru/asp/XML" + x.Name + ".asp"
		} else {
			if x.Name == "FindWMPurseNew" {
				url = "https://w10s.wmtransfer.com/asp/XMLFindWMPurseCertNew.asp"
			} else {
				url = "https://w10s.wmtransfer.com/asp/XML" + x.Name + "Cert.asp"
			}
		}
	} else if x.Type == "merchant" {
		url = "https://merchant.webmoney.ru/conf/xml/XML" + x.Name + ".asp"
	}
	return url
}

type requestW3s struct {
	XMLName xml.Name    `xml:"w3s.request"`
	Reqn    string      `xml:"reqn"`
	Wmid    string      `xml:"wmid"`
	Sign    string      `xml:"sign"`
	Request interface{} `xml:",>"`
}

//base struct for all response
type responseW3s struct {
	XMLName  xml.Name    `xml:"w3s.response"`
	Reqn     string      `xml:"reqn"`
	Retval   int64       `xml:"retval"`
	Retdesc  string      `xml:"retdesc"`
	Ser      string      `xml:"ser"`
	CWmid    string      `xml:"cwmid"`
	Wmid     string      `xml:"wmid"`
	Response interface{} `xml:",any"`
}
type merchantRequest struct {
	XMLName          xml.Name `xml:"merchant.request"`
	Wmid             string   `xml:"wmid"`
	Sign             string   `xml:"sign"`
	Sha256           string   `xml:"sha256"`
	Md5              string   `xml:"md5"`
	SecretKey        string   `xml:"secret_key"`
}

func (m *merchantRequest) setWmid(w string) {
	m.Wmid = w
}

func (m *merchantRequest) setSign(w string) {
	m.Sign = w
}

func (m *merchantRequest) setSha256(w string) {
	m.Sha256 = w
}

type merchantResponse struct {
	XmlName xml.Name `xml:"merchant.response"`
	ErrorLog ErrorLog    `xml:"errorlog"`
	Retval   int64       `xml:"retval"`
	Retdest  string      `xml:"retdesc"`
	Response interface{} `xml:",any"`
}

type ErrorLog struct {
	XmlName       xml.Name `xml:"errorlog"`
	LmiPayeePurse string   `xml:"lmi_payee_purse,attr"`
	LmiPaymentNo  string   `xml:"lmi_payment_no,attr"`
	Datecrt       string   `xml:"datecrt"`
	Dateupd       string   `xml:"dateupd"`
	DateS         string   `xml:"date_s"`
	DatePc        string   `xml:"date_pc"`
	DatePd        string   `xml:"date_pd"`
	PType         string   `xml:"p_type"`
	ErrCode       string   `xml:"err_code"`
	Siteid        string   `xml:"siteid"`
	Att           string   `xml:"att"`
	DateNotify    string   `xml:"date_notify"`
	ShopId        string   `xml:"shop_id"`
}
