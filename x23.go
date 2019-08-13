package webmoney

import (
	"encoding/xml"
)

type InvoiceRefuse struct {
	XMLName xml.Name `xml:"invoicerefuse`
	Wmid    string   `xml:"wmid"`
	Wminvid string   `xml:"wminvid"`
}

type InvoiceRefuseResponse struct {
	Id      string `xml:"id,attr"`
	Ts      string `xml:"ts,attr"`
	State   string `xml:"state"`
	DateUpd string `xml:dateupd`
}

func (w *WmClient) RefuseInvoice(i InvoiceRefuse) (InvoiceRefuseResponse, error) {
	w.Reqn = Reqn()
	w.X = X23
	if w.IsClassic() {
		w.Sign = i.Wmid + i.Wminvid + w.Reqn
	}
	w.Request = i
	result := InvoiceRefuseResponse{}
	err := w.getResult(&result)
	return result, err
}
