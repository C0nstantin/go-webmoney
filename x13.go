// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X13

package webmoney

import (
	"encoding/xml"
)

type RejectProtect struct {
	XMLName  xml.Name `xml:"rejectprotect"`
	WmTranId string   `xml:"wmtranid"`
}

/*
type Operation struct {
	XMLName  xml.Name `xml:"operation"`
	Id       string   `xml:"id, attr"`
	Ts       string   `xml:"ts,attr"`
	Opertype string   `xml:"opertype"`
	DateUpd  string   `xml:"dateupd"`
}
*/

func (w *WmClient) RejectProtect(i RejectProtect) (Operation, error) {
	w.Reqn = Reqn()
	w.X = X13
	if w.IsClassic() {
		w.Sign = i.WmTranId + w.Reqn
	}
	w.Request = i
	result := Operation{}
	err := w.getResult(&result)
	return result, err
}
