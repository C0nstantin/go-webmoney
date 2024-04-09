package tests

import "testing"

func TestClassicX2(t *testing.T) {
	wmClientClassic := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}

	transaction := webmoney.Transaction{
		TranId:    "1001",
		PurseDest: "R305844543288",
		PurseSrc:  "R562831674289",
		Amount:    "3.2",
		Period:    "1",
		Desc:      "Теst кирилица",
		PCode:     "",
		WmInvid:   "0",
		OnlyAuth:  "1",
	}

	result, err := wmClientClassic.CreateTransaction(transaction)

	t.Log(result)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLightX2(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Key:  CnfLight.Key,
		Cert: CnfLight.Cert,
	}

	transaction := webmoney.Transaction{
		TranId:    "1002",
		PurseDest: "Z311164609945",
		PurseSrc:  "Z214605808008",
		Amount:    "1.03",
		Period:    "1",
		Desc:      "Текст кирилица",
		Pcode:     "",
		WmInvid:   "0",
		OnlyAuth:  "1",
	}

	result, err := wmCl.CreateTransaction(transaction)

	t.Log(result)
	if err != nil {
		t.Fatal(err)
	}
}
