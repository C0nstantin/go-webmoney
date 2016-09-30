// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X9

package webmoney

import (
	"encoding/xml"
)

type Purses struct {
	XMLName xml.Name `xml:"getpurses"`
	Wmid    string   `xml:"wmid"`
}

type PursesResp struct {
	XMLName   xml.Name `xml:"purses"`
	Cnt       string   `xml:"cnt,attr"`
	PurseList []RPurse `xml:"purses"`
}

type RPurse struct {
	XMLName     xml.Name `xml:"purse"`
	PurseName   string   `xml:"pursename"`
	Amount      string   `xml:"amount"`
	Desc        string   `xml:"desc"`
	Outsideopen string   `xml:"outsideopen"`
	LastIntr    string   `xml:"lastintr"`
	LastOuttr   string   `xml:"lastouttr"`
}

func (w *WmClient) GetPurses4Wmid(p Purses) (PursesResp, error) {

	w.Reqn = Reqn()
	w.X = X9
	if w.IsClassic() {
		w.Sign = p.Wmid + w.Reqn
	}
	w.Request = p
	result := PursesResp{}
	err := w.getResult(&result)
	return result, err

}
