package webmoney

import (
	"encoding/xml"
)

type FinishProtect struct {
	XMLName  xml.Name `xml:"finishprotect"`
	WmTranId string   `xml:"wmtranid"`
	PCode    string   `xml:"pcode"`
}

/*
type Operation struct {
	XMLName  xml.Name `xml:operation`
	Id       string   `xml:"id,attr"`
	Ts       string   `xml:"ts,attr"`
	OperType string   `xml:"opertype"`
	DateUpd  string   `xml:"dateupd"`
}
*/
func (w *WmClient) DoFinishProtect(f FinishProtect) (Operation, error) {
	w.Reqn = Reqn()
	w.X = X5

	if w.IsClassic() {
		w.Sign = f.PCode + f.WmTranId + w.Reqn
	}

	w.Request = f
	result := Operation{}
	err := w.getResult(&result)
	return result, err
}
