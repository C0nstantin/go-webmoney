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

func (i RejectProtect) GetSignSource(reqn string) (string, error) {
	return i.WmTranId + reqn, nil
}

func (w *WmClient) RejectProtect(i RejectProtect) (Operation, error) {

	X := W3s{
		Interface: XInterface{Name: "RejectProtect", Type: "w3s"},
		Request:   i,
		Client:    w,
	}

	result := Operation{}
	err := X.getResult(&result)
	return result, err
}
