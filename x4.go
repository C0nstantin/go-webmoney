// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// Receiving the history of issued invoices. Verifying whether invoices were paid.
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X4

package webmoney

import (
	"encoding/xml"
)

type OutInvoices struct {
	XMLName    xml.Name `xml:"getoutinvoices"`
	Purse      string   `xml:"purse"`
	WmInvid    string   `xml:"wminvid"`
	OrderId    string   `xml:"orderid"`
	DateStart  string   `xml:"datestart"`
	DateFinish string   `xml:"datefinish"`
}

type OutInvoicesResp struct {
	XMLName     xml.Name          `xml:"outinvoices"`
	InvoiceList []InvoiceResponse `xml:"outinvoice"`
}

//type OutInvoice struct {
//	XMLName xml.Name `xml:outinvoice`
//	InvoiceResponse
//}

func (w *WmClient) GetOutInvoices(i OutInvoices) (OutInvoicesResp, error) {
	w.Reqn = Reqn()
	w.X = X4
	if w.IsClassic() {
		w.Sign = i.Purse + w.Reqn
	}
	w.Request = i
	result := OutInvoicesResp{}
	err := w.getResult(&result)
	return result, err
}
