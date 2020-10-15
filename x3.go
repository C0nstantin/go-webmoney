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

func (o GetOpers) GetSignSource(reqn string) (string, error) {
	return o.Purse + reqn, nil
}

type Operations struct {
	XMLName       xml.Name    `xml:"operations"`
	OperationList []Operation `xml:"operation"`
}

func (w *WmClient) GetOperations(o GetOpers) (Operations, error) {

	X := W3s{
		Interface: XInterface{Name: "Operations", Type: "w3s"},
		Request:   o,
		Client:    w,
	}

	result := Operations{}

	err := X.getResult(&result)
	return result, err
}
