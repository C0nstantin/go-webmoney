// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// Receiving Transaction History. Checking Transaction Status.
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X3

package webmoney

import (
	"encoding/xml"
)

type GetOpers struct {
	XMLName    xml.Name `xml:"getoperations"`
	Purse      string   `xml:"purse"`
	WmTranId   string   `xml:"wmtranid"`
	TranId     string   `xml:"tranid"`
	WmInvid    string   `xml:"wminvid"`
	OrderId    string   `xml:"orderid"`
	DateStart  string   `xml:"datestart"`
	DateFinish string   `xml:"datefinish"`
}

type Operations struct {
	XMLName       xml.Name    `xml:"operations"`
	OperationList []Operation `xml:"operation"`
}

func (w *WmClient) GetOperations(o GetOpers) (Operations, error) {
	w.Reqn = Reqn()
	w.X = X3
	if w.IsClassic() {
		w.Sign = o.Purse +
			w.Reqn
	}
	w.Request = o
	result := Operations{}

	err := w.getResult(&result)

	return result, err
}
