package webmoney

import (
	"encoding/xml"
)

type FinishProtect struct {
	XMLName  xml.Name `xml:"finishprotect"`
	WmTranId string   `xml:"wmtranid"`
	PCode    string   `xml:"pcode"`
}

func (f FinishProtect) GetSignSource(reqn string) (string, error) {
	return f.WmTranId + f.PCode + reqn, nil
}

func (w *WmClient) DoFinishProtect(f FinishProtect) (Operation, error) {

	X := W3s{
		Request:   f,
		Interface: XInterface{Name: "FinishProtect", Type: "w3s"},
		Client:    w,
	}

	result := Operation{}
	err := X.getResult(&result)
	return result, err
}
