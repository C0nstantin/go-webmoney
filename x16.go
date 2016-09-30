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

type PurseResponse struct {
	XMLName   xml.Name `xml:"purse"`
	Id        string   `xml:"id,attr"`
	PurseName string   `xml:"pursename"`
	Amount    string   `xml:"amount"`
	Desc      string   `xml:"desc"`
}

func (w *WmClient) CreatePurse(p Purse) (PurseResponse, error) {
	w.Reqn = Reqn()
	w.X = X16
	if w.IsClassic() {
		w.Sign = p.Wmid + p.PurseType + w.Reqn
	}
	w.Request = p
	result := PurseResponse{}
	err := w.getResult(&result)
	return result, err
}
