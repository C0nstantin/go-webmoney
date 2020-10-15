package webmoney

import (
	"encoding/xml"
	"strings"
)

type SendMsg struct {
	XMLName      xml.Name `xml:"message"`
	ReceiverWmid string   `xml:"receiverwmid"`
	MsgSubj      string   `xml:"msgsubj"`
	MsgText      string   `xml:"msgtext"`
	OnlyAuth     int      `xml:"onlyauth"`
}

func (s SendMsg) GetSignSource(reqn string) (string, error) {
	subject, err := Utf8ToWin(strings.Trim(s.MsgSubj, " "))
	if err != nil {
		return "", err
	}

	text, err := Utf8ToWin(strings.Trim(s.MsgText, " "))
	if err != nil {
		return "", err
	}

	return s.ReceiverWmid + reqn + text + subject, nil

}

type SendMsgResponse struct {
	XMLName      xml.Name `xml:"message"`
	MessageId    int      `xml:"id,attr"`
	ReceiverWmid string   `xml:"receiverwmid"`
	MsgSubj      string   `xml:"msgsubj"`
	MsgText      string   `xml:"msgtext"`
	DateCrt      string   `xml:"datecrt"`
}

func (w *WmClient) SendMessage(s SendMsg) (SendMsgResponse, error) {

	X := W3s{
		Request:   s,
		Interface: XInterface{Name: "SendMsg", Type: "w3s"},
		Client:    w,
	}

	result := SendMsgResponse{}
	err := X.getResult(&result)
	return result, err
}
