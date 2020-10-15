// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X16

package webmoney

import (
	"encoding/xml"
)

type Purse struct {
	XMLName   xml.Name `xml:"createpurse"`
	Wmid      string   `xml:"wmid"`
	PurseType string   `xml:"pursetype"`
	Desc      string   `xml:"desc"`
}

func (p Purse) GetSignSource(reqn string) (string, error) {
	return p.Wmid + p.PurseType + reqn, nil
}

type PurseResponse struct {
	XMLName   xml.Name `xml:"purse"`
	Id        string   `xml:"id,attr"`
	PurseName string   `xml:"pursename"`
	Amount    string   `xml:"amount"`
	Desc      string   `xml:"desc"`
}

func (w *WmClient) CreatePurse(p Purse) (PurseResponse, error) {
	X := W3s{
		Interface: XInterface{Name: "CreatePurse", Type: "w3s"},
		Client:    w,
		Request:   p,
	}

	result := PurseResponse{}
	err := X.getResult(&result)
	return result, err
}
