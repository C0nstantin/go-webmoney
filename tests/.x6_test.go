package tests

import "testing"
import "webmoney"
import "fmt"

func TestX6(t *testing.T) {

	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}

	wmMsg := webmoney.SendMsg{
		ReceiverWmid: "128756507061",
		MsgSubj:      "testClassic",
		MsgText:      "Тест сообщения через Классик кирилица",
	}

	result, err := wmCl.SendMessage(wmMsg)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)

	if result.MessageId <= 0 {
		t.Log(result)
		t.Fatal("Classic fatal test Message ID is empty !!! ")
	}

	wmClight := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Cert: CnfLight.Cert,
		Key:  CnfLight.Key,
	}

	wmMsg1 := webmoney.SendMsg{
		ReceiverWmid: "128756507061",
		MsgSubj:      "testLight",
		MsgText:      "Тест сообщения через Light кирилица",
	}

	result1, err := wmClight.SendMessage(wmMsg1)

	fmt.Println(wmClight.ResultStr)

	if err != nil {
		t.Fatal(err)
	}

	if result1.MessageId <= 0 {
		t.Log(result)
		t.Fatal("Light fatal Message ID is empty !!! ")
	}

	fmt.Println(result)
}
