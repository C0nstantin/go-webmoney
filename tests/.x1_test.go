//finish
package tests

import "testing"
import "webmoney"
import "fmt"
import "strconv"

func TestX1(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	invoice := webmoney.Invoice{
		OrderId:      "10",
		CustomerWmid: "128756507061",
		StorePurse:   "R562831674289",
		Amount:       "3.2",
		Desc:         "Тест кирилица",
		Address:      "Россия, Москва, NJ переулок д 1",
		Period:       "1",
		Expiration:   "1",
		OnlyAuth:     "0",
	}
	result, err := wmCl.SendInvoice(invoice)
	fmt.Println(result)
	t.Log(result)
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := strconv.ParseInt(result.Id, 10, 64); id <= 0 {
		t.Log(result)
		t.Fatal(" Light Fail Invoice ID is empty !!! ")
	}
	wmClight := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Key:  CnfLight.Key,
		Cert: CnfLight.Cert,
	}
	invoicel := webmoney.Invoice{
		OrderId:      "10",
		CustomerWmid: "128756507061",
		StorePurse:   "Z214605808008",
		Amount:       "3.2",
		Desc:         "Теcт кирилица light",
		Address:      "Россия, Москва, ТТ переулок д 1",
		Period:       "1",
		Expiration:   "1",
		OnlyAuth:     "0",
	}
	resultl, err := wmClight.SendInvoice(invoicel)
	fmt.Println(resultl)
	t.Log(resultl)
	fmt.Println(wmClight.ResultStr)
	if err != nil {
		t.Fatal(err)
	}
	if id, _ := strconv.ParseInt(resultl.Id, 10, 64); id <= 0 {
		t.Log(result)
		t.Fatal("Light Fail Invoice ID is empty !!! ")
	}
}
