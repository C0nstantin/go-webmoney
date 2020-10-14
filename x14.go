// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// Sending invoice from merchant to customer.
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X14

package webmoney

import (
	"encoding/xml"
)

type Trans struct {
	XMLName        xml.Name `xml:"trans"`
	InWmTranId     string   `xml:"inwmtranid"`
	Amount         string   `xml:"amount"`
	MoneyBackPhone string   `xml:"moneybackphone"`
}

func (t Trans) GetSignSource(reqn string) (string, error) {
	return reqn + t.InWmTranId + t.Amount, nil
}

func (w *WmClient) MoneyBack(t Trans) (Operation, error) {

	X := W3s{
		Interface: XInterface{Name: "TransMoneyback", Type: "w3s"},
		Request:   t,
		Client:    w,
	}

	result := Operation{}

	err := X.getResult(&result)
	return result, err
}
