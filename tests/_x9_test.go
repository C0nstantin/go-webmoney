package tests

import "testing"
import "webmoney"

func TestX9(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	p := webmoney.Purses{
		Wmid: CnfClassic.Wmid,
	}
	result, err := wmCl.GetPurses4Wmid(p)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
