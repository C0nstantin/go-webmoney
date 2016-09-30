// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// Sending invoice from merchant to customer.
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X1

package webmoney

import (
	"encoding/xml"
)

type Invoice struct {
	XMLName      xml.Name `xml:"invoice"`
	OrderId      string   `xml:"orderid"`
	CustomerWmid string   `xml:"customerwmid"`
	StorePurse   string   `xml:"storepurse"`
	Amount       string   `xml:"amount"`
	Desc         string   `xml:"desc"`
	Address      string   `xml:"address"`
	Period       string   `xml:"period"`
	Expiration   string   `xml:"expiration"`
	OnlyAuth     string   `xml:"onlyauth"`
	Lmi_shop_id  string   `xml:"lmi_shop_id"`
}

type InvoiceResponse struct {
	Id           string `xml:"id,attr"`
	Ts           string `xml:"ts,attr"`
	OrderId      string `xml:"orderid"`
	CustomerWmid string `xml:"customerwmid"`
	StorePurse   string `xml:"storepurse"`
	Amount       string `xml:"amount"`
	Desc         string `xml:"desc"`
	Address      string `xml:"address"`
	Period       string `xml:"period"`
	Expiration   string `xml:"expiration"`
	State        string `xml:"state"`
	DateCrt      string `xml:"datecrt"`
	DateUpd      string `xml:"dateupd"`
	WmTranId     string `xml:"wmtranid"`
}

func (w *WmClient) SendInvoice(i Invoice) (InvoiceResponse, error) {
	w.Reqn = Reqn()
	w.X = X1
	if w.IsClassic() {
		//orderid+customerwmid+storepurse+amount+desc+address+period+expiration+reqn
		desc, err := Utf8ToWin(i.Desc)
		if err != nil {
			return InvoiceResponse{}, err
		}
		address, err := Utf8ToWin(i.Address)
		if err != nil {
			return InvoiceResponse{}, err
		}

		w.Sign = string(i.OrderId) +
			i.CustomerWmid +
			i.StorePurse +
			i.Amount +
			desc +
			address +
			i.Period +
			i.Expiration +
			w.Reqn
	}
	w.Request = i
	result := InvoiceResponse{}

	err := w.getResult(&result)
	return result, err
}
