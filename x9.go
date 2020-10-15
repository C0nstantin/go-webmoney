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

func (p Purses) GetSignSource(reqn string) (string, error) {

	return p.Wmid + reqn, nil
}

type PursesResp struct {
	XMLName   xml.Name `xml:"purses"`
	Cnt       string   `xml:"cnt,attr"`
	PurseList []RPurse `xml:"purse"`
}

type RPurse struct {
	XMLName          xml.Name `xml:"purse"`
	Id               string   `xml:"id,attr"`
	PurseName        string   `xml:"pursename"`
	Amount           string   `xml:"amount"`
	Desc             string   `xml:"desc"`
	Outsideopen      string   `xml:"outsideopen"`
	Outsideopenstate string   `xml:"outsideopenstate"`
	Byptlimit        string   `xml:"byptlimit"`
	LastIntr         string   `xml:"lastintr"`
	LastOuttr        string   `xml:"lastouttr"`
}

func (w *WmClient) GetPurses4Wmid(p Purses) (PursesResp, error) {

	X := W3s{
		Interface: XInterface{Name: "Purses", Type: "w3s"},
		Request:   p,
		Client:    w,
	}

	result := PursesResp{}

	err := X.getResult(&result)

	return result, err
}
