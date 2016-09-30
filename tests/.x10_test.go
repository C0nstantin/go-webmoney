package tests

import "testing"
import "webmoney"

func TestX9(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Cert: CnfLight.Cert,
		Key:  CnfLight.Key,
	}
	i := webmoney.InInvoices{
		Wmid:       CnfLight.Wmid,
		WmInvid:    "0",
		DateStart:  "20140101 00:00:00",
		DateFinish: "20160930 00:00:00",
	}
	result, err := wmCl.GetInInvoices(i)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}
