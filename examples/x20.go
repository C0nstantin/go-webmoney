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
	got := webmoney.X20Request{
		LmiPayeePurse:        "Z303339773989",
		LmiPaymentNo:         "121",
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
	}

	result, err := wmClient.TransRequest(got)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
	fmt.Println("Insert code for confirmation")
	var code string
	fmt.Scan(&code)
	got1 := webmoney.X202Request{
		LmiPayeePurse:       "Z303339773989",
		LmiClientnumberCode: code,
		LmiWminvoiceid:      result.Wminvoiceid,
	}
	result1, err := wmClient.TransConfirm(got1)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result1)
}
