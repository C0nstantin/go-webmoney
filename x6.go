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
type SendMsgResponse struct {
	XMLName      xml.Name `xml:"message"`
	MessageId    int      `xml:"id,attr"`
	ReceiverWmid string   `xml:"receiverwmid"`
	MsgSubj      string   `xml:"msgsubj"`
	MsgText      string   `xml:"msgtext"`
	DateCrt      string   `xml:"datecrt"`
}

func (w *WmClient) SendMessage(s SendMsg) (SendMsgResponse, error) {
	w.Reqn = Reqn() //in request metrhod
	w.X = X6
	subject := strings.Trim(s.MsgSubj, " ")
	text := strings.Trim(s.MsgText, " ")

	if w.IsClassic() {
		subject, _ = Utf8ToWin(subject)
		text, _ = Utf8ToWin(text)
		w.Sign = s.ReceiverWmid + w.Reqn + text + subject
		text, _ = WinToUtf8(text)
		subject, _ = WinToUtf8(subject)
	}

	w.Request = s
	result := SendMsgResponse{}

	err := w.getResult(&result)
	return result, err
}
