// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://w3s.webmoney.ru/asp/XMLGetCorrList.asp

package webmoney

import (
	"encoding/xml"
	"errors"
)

type CorrList struct {
	XMLName xml.Name `xml:"getcorrlist"`
	Wmid    string   `xml:"wmid"`
}

type CorrListResponse struct {
	XMLName xml.Name `xml:"corrlist"`
	Count   int      `xml:"cnt,attr"`
	Corrs   []Corr   `xml:"corr"`
}

type Corr struct {
	Wmid string `xml:"wmid"`
	Nick string `xml:"nick"`
}

func (c CorrList) GetSignSource(reqn string) (string, error) {
	return c.Wmid + reqn, nil
}

func (w *WmClient) GetCorrList(c CorrList) (CorrListResponse, error) {
	result := CorrListResponse{}
	if !w.IsClassic() {
		return result, errors.New("not support Keeper Web Pro")
	}

	X := W3s{
		Interface: XInterface{Name: "GetCorrList", Type: "w3s"},
		Request:   c,
		Client:    w,
	}

	err := X.getResult(&result)
	return result, err
}
