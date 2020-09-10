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

type InInvoicesResponse struct {
	XMLName     xml.Name          `xml:"ininvoices"`
	InvoiceList []InvoiceResponse `xml:ininvoices`
}

func (w *WmClient) GetInInvoices(i InInvoices) (InInvoicesResponse, error) {
	w.Reqn = Reqn()
	w.X = X10
	if w.IsClassic() {
		w.Sign = i.Wmid + i.WmInvid + i.DateStart + i.DateFinish + w.Reqn
	}

	w.Request = i
	result := InInvoicesResponse{}
	err := w.getResult(&result)
	return result, err
}
