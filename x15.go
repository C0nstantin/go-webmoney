// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X15

package webmoney

import (
	"encoding/xml"
)

type TrustList struct {
	XMLName xml.Name `xml:"gettrustlist"`
	Wmid    string   `xml:"wmid"`
}

func (t TrustList) GetSignSource(reqn string) (string, error) {
	return t.Wmid + reqn, nil
}

type TrustListResponse struct {
	XMLName   xml.Name `xml:"trustlist"`
	Cnt       string   `xml:"cnt,attr"`
	TrustList []Trust  `xml:"trastlist"`
}

type Trust struct {
	XMLName     xml.Name `xml:"trast"`
	Id          string   `xml:"id,attr"`
	Inv         string   `xml:"inv,attr"`
	Trast       string   `xml:"trans,attr"`
	PurseAttr   string   `xml:"purse,attr"`
	TransHist   string   `xml:"transhist"`
	Master      string   `xml:"master"`
	Purse       string   `xml:"purse"`
	DayLimit    string   `xml:"daylimit"`
	DLimit      string   `xml:"dlimit"`
	WLimit      string   `xml:"wlimit"`
	MLimit      string   `xml:"mlimit"`
	DSum        string   `xml:"dsum"`
	WSum        string   `xml:"wsum"`
	MSum        string   `xml:"msum"`
	LastSumDate string   `xml:"lastsumdate"`
}

type TrustSave struct {
	XMLName    xml.Name `xml:"trust"`
	Inv        string   `xml:"inv,attr"`
	Trans      string   `xml:"trans,attr"`
	PurseAttr  string   `xml:"purse,attr"`
	TransHist  string   `xml:"transhist"`
	MasterWmid string   `xml:"masterwmid"`
	SlaveWmid  string   `xml:"slavewmid"`
	Purse      string   `xml:"purse"`
	Limit      string   `xml:"limit"`
	DayLimit   string   `xml:"daylimit"`
	WeekLimit  string   `xml:"weeklimit"`
	MonthLimit string   `xml:"monthlimit"`
}

func (t TrustSave) GetSignSource(reqn string) (string, error) {
	return t.Purse + t.MasterWmid + reqn, nil
}

//func GetTrastList
func (w *WmClient) GetTrustList(t TrustList) (TrustListResponse, error) {
	var InterfaceName string

	if w.Wmid != t.Wmid {
		InterfaceName = "TrustList2"
	} else {
		InterfaceName = "TrustList"
	}
	X := W3s{
		Interface: XInterface{Name: InterfaceName, Type: "w3s"},
		Request:   t,
		Client:    w,
	}
	result := TrustListResponse{}
	err := X.getResult(&result)
	return result, err
}

func (w *WmClient) SetTrust(t TrustSave) (Trust, error) {
	X := W3s{
		Interface: XInterface{Name: "TrustSave2", Type: "w3s"},
		Request:   t,
		Client:    w,
	}

	result := Trust{}
	err := X.getResult(&result)
	return result, err
}
