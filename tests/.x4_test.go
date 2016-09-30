package tests

import "testing"
import "webmoney"
import "fmt"

func TestX4(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	invoices := webmoney.OutInvoices{
		Purse:      "R562831674289",
		DateStart:  "20160622 00:00:00",
		DateFinish: "20160927 00:00:00",
	}
	result, err := wmCl.GetOutInvoices(invoices)

	t.Log(result)
	fmt.Println(result)
	fmt.Println(wmCl.ResultStr)
	if err != nil {
		t.Fatal(err)
	}

}
