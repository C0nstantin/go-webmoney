package tests

import "testing"
import "webmoney"
import "fmt"

func TestX3(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	opts := webmoney.GetOpers{
		TranId:     "10",
		Purse:      "R562831674289",
		DateStart:  "20160726 00:00:00",
		DateFinish: "20160721 00:00:00",
	}
	result, err := wmCl.GetOperations(opts)
	fmt.Println(wmCl.ResultStr)

	t.Log(result)
	if err != nil {
		t.Fatal(err)
	}

	wmClight := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Cert: CnfLight.Cert,
		Key:  CnfLight.Key,
	}
	transaction1 := webmoney.GetOpers{
		Purse:      "Z311164609945",
		DateStart:  "20160726 00:00:00",
		DateFinish: "20160922 00:00:00",
	}
	result1, err := wmClight.GetOperations(transaction1)
	fmt.Println(wmClight.ResultStr)
	fmt.Println(result1)
	if err != nil {
		t.Fatal(err)
	}

}
