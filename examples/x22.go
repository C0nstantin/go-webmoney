package main

import (
	"fmt"
	"github.com/C0nstantin/go-webmoney"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"log"
)

func main() {
	date, err := ioutil.ReadFile("examples/conf.toml")
	if err != nil {
		log.Fatal(err)
	}

	config, _ := toml.Load(string(date))
	wmClient := webmoney.WmClient{
		Wmid: config.Get("client1.Wmid").(string),
		Key:  config.Get("client1.Key").(string),
		Pass: config.Get("client1.Pass").(string),
	}
	got := webmoney.X22Request{

		Validityperiodinhours: "1",
		Paymenttags: webmoney.Paymenttags{

			LmiPayeePurse:        "Z303339773989",
			LmiPaymentNo:         "122",
			LmiPaymentNoType:     "",
			LmiPaymentAmount:     "0.1",
			LmiPaymentDesc:       "test test test",
			LmiPaymentDescBase64: "",
			LmiClientnumber:      "79261919656",
			LmiClientnubmerType:  "0",
			LmiSmsType:           "3",
			LmiShopId:            "",
			LmiHold:              "",
			Lang:                 "",
			EnulatedFlag:         "",
		},
	}

	result, err := wmClient.TransSave(got)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
	fmt.Println("https://merchant.wmtransfer.com/lmi/payment.asp?gid=" + result.Transtoken)

}
