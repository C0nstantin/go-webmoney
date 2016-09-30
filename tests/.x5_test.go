package tests

import "testing"

import "webmoney"

func TestX5(t *testing.T) {

	//create transactionf with pcode

	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	pcode := "testtest"
	tt := webmoney.Transaction{
		TranId:    "1999",
		PurseDest: "R988180789538",
		PurseSrc:  "R562831674289",
		Amount:    "3.2",
		Period:    "1",
		Desc:      "Теst кирилица",
		WmInvid:   "0",
		OnlyAuth:  "1",
		PCode:     pcode,
	}

	result, err := wmCl.CreateTransaction(tt)

	if err != nil {
		t.Fatal(err)
	}

	wmClight := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Cert: CnfLight.Cert,
		Key:  CnfLight.Key,
	}

	fp := webmoney.FinishProtect{
		WmTranId: result.Ts,
		PCode:    pcode,
	}

	result1, err := wmClight.DoFinishProtect(fp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result1)

}
