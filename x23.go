package webmoney

import (
	"encoding/xml"
)

type InvoiceRefuse struct {
	XMLName xml.Name `xml:"invoicerefuse"`
	Wmid    string   `xml:"wmid"`
	Wminvid string   `xml:"wminvid"`
}

func (i InvoiceRefuse) GetSignSource(reqn string) (string, error) {
	return i.Wmid + i.Wminvid + reqn, nil
}

type InvoiceRefuseResponse struct {
	Id      string `xml:"id,attr"`
	Ts      string `xml:"ts,attr"`
	State   string `xml:"state"`
	DateUpd string `xml:"dateupd"`
}

func (w *WmClient) RefuseInvoice(i InvoiceRefuse) (InvoiceRefuseResponse, error) {
	X := W3s{
		Client:    w,
		Interface: XInterface{Name: "InvoiceRefusal", Type: "w3s"},
		Request:   i,
	}
	result := InvoiceRefuseResponse{}
	err := X.getResult(&result)
	return result, err
}
