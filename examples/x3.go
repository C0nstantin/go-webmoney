//Webmoney XML Interfaces
//Interface X1. Sending invoice from merchant to customer.

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
	got := webmoney.GetOpers{
		Purse:      "Z303339773000",
		WmTranId:   "",
		TranId:     "",
		WmInvid:    "",
		OrderId:    "",
		DateStart:  "20200911",
		DateFinish: "20201013",
	} //1766082859

	result, err := wmClient.GetOperations(got)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
}
