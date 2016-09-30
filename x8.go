package webmoney

// Copyright (c) 2016 The Constantin Karataev. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd.
//
// desc
// https://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X8

/*<w3s.request>
    <reqn></reqn>
    <wmid></wmid>
    <sign></sign>
    <testwmpurse>
        <wmid></wmid>
        <purse></purse>
    </testwmpurse>
</w3s.request>*/
/*
<w3s.response>
    <reqn></reqn>
        <retval></retval>
        <retdesc></retdesc>
        <testwmpurse>
            <wmid available="-1" themselfcorrstate="-1" newattst="-1"> </wmid>
            <purse merchant_active_mode="-1" merchant_allow_cashier="-1"></purse>
        </testwmpurse>
</w3s.response>
*/
import (
	"encoding/xml"
)

// request
type TestWmPurse struct {
	XMLName xml.Name `xml:"testwmpurse"`
	Wmid    string   `xml:"wmid"`
	Purse   string   `xml:"purse"`
}

type TestWmPurseResponse struct {
	XMLName xml.Name      `xml:"testwmpurse"`
	Wmid    Wmid          `xml:"wmid"`
	Purse   ReturnedPurse `xml:"purse"`
}
type ReturnedPurse struct {
	XMLName              xml.Name `xml:"purse"`
	Value                string   `xml:",chardata"`
	MerchantActiveMode   string   `xml:"merchant_active_mode,attr"`
	MerchantAllowCashier string   `xml:"merchant_allow_cashier,attr"`
}
type Wmid struct {
	XMLName           xml.Name `xml:"wmid"`
	Value             string   `xml:",chardata"`
	Available         string   `xml:"available,attr"`
	Themselfcorrstate string   `xml:"themselfcorrstate,attr"`
	Newattst          string   `xml:"newattst,attr"`
}

func (w *WmClient) FindWmidPurse(t TestWmPurse) (TestWmPurseResponse, error) {
	w.X = X8
	w.Reqn = Reqn()
	if w.IsClassic() {
		w.Sign = t.Wmid + t.Purse
	}
	w.Request = t
	result := TestWmPurseResponse{}
	err := w.getResult(&result)
	return result, err
}
