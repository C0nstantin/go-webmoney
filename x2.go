// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.

// Create Transaction
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X2

package webmoney

import (
	"encoding/xml"
)

type Transaction struct {
	XMLName   xml.Name `xml:"trans"`
	TranId    string   `xml:"tranid"`
	PurseSrc  string   `xml:"pursesrc"`
	PurseDest string   `xml:"pursedest"`
	Amount    string   `xml:"amount"`
	Period    string   `xml:"period"`
	Desc      string   `xml:"desc"`
	PCode     string   `xml:"pcode"`
	WmInvid   string   `xml:"wminvid"`
	OnlyAuth  string   `xml:"onlyauth"`
}

type Operation struct {
	XMLName   xml.Name `xml:operation`
	Id        string   `xml:"id,attr"`
	Ts        string   `xml:"ts,attr"`
	TranId    string   `xml:"tranid"`
	PurseSrc  string   `xml:"pursesrc"`
	PurseDest string   `xml:"pursedest"`
	Amount    string   `xml:"amount"`
	Commis    string   `xml:"comiss"`
	Opertype  string   `xml:"opertype"`
	Period    string   `xml:"period"`
	WmInvid   string   `xml:"wminvid"`
	Desc      string   `xml:"desc"`
	DateCrt   string   `xml:"datecrt"`
	DateUpd   string   `xml:"dateupd"`
	CorrWm    string   `xml:"corrwm"`
	Rest      string   `xml:"rest"`
	TimeLock  bool     `xml:timelock`
}

func (w *WmClient) CreateTransaction(t Transaction) (Operation, error) {
	w.Reqn = Reqn()
	w.X = X2

	if w.IsClassic() {
		desc, err := Utf8ToWin(t.Desc)
		if err != nil {
			return Operation{}, err
		}
		w.Sign = w.Reqn +
			t.TranId +
			t.PurseSrc +
			t.PurseDest +
			t.Amount +
			t.Period +
			t.PCode +
			desc +
			t.WmInvid
	}

	w.Request = t
	result := Operation{}
	err := w.getResult(&result)
	return result, err

}
