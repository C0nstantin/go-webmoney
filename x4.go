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

func (i OutInvoices) GetSignSource(reqn string) (string, error) {
	return i.Purse + reqn, nil
}

type OutInvoicesResp struct {
	XMLName     xml.Name          `xml:"outinvoices"`
	InvoiceList []InvoiceResponse `xml:"outinvoice"`
}

func (w *WmClient) GetOutInvoices(i OutInvoices) (OutInvoicesResp, error) {

	X := W3s{
		Request:   i,
		Interface: XInterface{Name: "OutInvoices", Type: "w3s"},
		Client:    w,
	}

	result := OutInvoicesResp{}
	err := X.getResult(&result)
	return result, err
}
