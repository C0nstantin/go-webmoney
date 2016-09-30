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
	XMName         xml.Name `xml:"trans"`
	InWmTranId     string   `xml:"inwmtranid"`
	Amount         string   `xml:"amount"`
	MoneyBackPhone string   `xml:"moneybackphone"`
}

func (w *WmClient) MoneyBack(t Trans) (Operation, error) {
	w.Reqn = Reqn()
	w.X = X14

	if w.IsClassic() {
		w.Sign = w.Reqn + t.InWmTranId + t.Amount
	}

	w.Request = t
	result := Operation{}

	err := w.getResult(&result)
	return result, err
}
