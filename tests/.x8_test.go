package tests

import "testing"
import "webmoney"

import "fmt"

func TestX8(t *testing.T) {
	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	testWm := webmoney.TestWmPurse{
		Wmid: `128756507061`,
	}

	testPurse := webmoney.TestWmPurse{
		Purse: `Z214605808008`,
	}
	testWmidPurse := webmoney.TestWmPurse{
		Wmid:  `128756507061`,
		Purse: `K133612664763`,
	}

	result, err := wmCl.FindWmidPurse(testWm)
	fmt.Println(result)
	t.Log(result)
	if err != nil {
		t.Fatal(err)
	}

	resultPurse, err1 := wmCl.FindWmidPurse(testPurse)
	fmt.Println(resultPurse)
	if err1 != nil {
		t.Fatal(err1)
	}
	if resultPurse.Purse.Value == "" {
		t.Fatal("Purse not valid")
	} else {
		fmt.Println(resultPurse.Purse.Value)
	}
	resultWmPurse, err2 := wmCl.FindWmidPurse(testWmidPurse)
	fmt.Println(resultWmPurse)
	if err2 != nil {
		t.Fatal(err2)
	}
	if resultWmPurse.Wmid.Value != `128756507061` || resultWmPurse.Purse.Value != `K133612664763` {
		t.Fatal("test with wmid and purse not valid")
	} else {
		fmt.Println(resultWmPurse.Wmid.Value)
	}
}
