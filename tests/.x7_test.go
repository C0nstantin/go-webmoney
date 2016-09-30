package tests

import "testing"
import "webmoney"
import "webmoney/wmsigner"

func TestX7(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	signer := wmsigner.NewSigner(CnfClassic.Wmid, CnfClassic.Wmid_pass, CnfClassic.Wmid_key)
	teststring, err := signer.Sign("R562831674289")
	if err != nil {
		t.Fatal(err)
	}
	signature := webmoney.TestSignRequest{
		Wmid: "128756507061",
		Plan: "R562831674289",
		Sign: teststring,
	}
	result, err := wmCl.TestSign(signature)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(result)
	}
}
