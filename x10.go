// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X10

package webmoney

import (
	"encoding/xml"
)

type InInvoices struct {
	XMLName    xml.Name `xml:"getininvoices"`
	Wmid       string   `xml:"wmid"`
	WmInvid    string   `xml:"wminvid"`
	DateStart  string   `xml:"datestart"`
	DateFinish string   `xml:"datefinish"`
}

func (i InInvoices) GetSignSource(reqn string) (string, error) {
	return i.Wmid + i.WmInvid + i.DateStart + i.DateFinish + reqn, nil
}

type InInvoicesResponse struct {
	XMLName     xml.Name          `xml:"ininvoices"`
	InvoiceList []InvoiceResponse `xml:"ininvoices"`
}

func (w *WmClient) GetInInvoices(i InInvoices) (InInvoicesResponse, error) {
	X := W3s{
		Interface: XInterface{Name: "InInvoices", Type: "w3s"},
		Request:   i,
		Client:    w,
	}

	result := InInvoicesResponse{}
	err := X.getResult(&result)
	return result, err
}
