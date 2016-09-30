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
	LastSumDate string   `xml:lastsumdate`
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

//func GetTrastList
func (w *WmClient) GetTrustList(t TrustList) (TrustListResponse, error) {
	w.Reqn = Reqn()
	if w.Wmid != t.Wmid {
		w.X = X152
	} else {
		w.X = X15
	}
	if w.IsClassic() {
		w.Sign = t.Wmid + w.Reqn
	}
	w.Request = t
	result := TrustListResponse{}
	err := w.getResult(&result)
	return result, err

}

func (w *WmClient) SetTrust(t TrustSave) (Trust, error) {
	w.Reqn = Reqn()
	w.X = X153
	if w.IsClassic() {
		w.Sign = w.Wmid + t.Purse + t.MasterWmid + w.Reqn
	}

	w.Request = t
	result := Trust{}
	err := w.getResult(&result)
	return result, err

}

