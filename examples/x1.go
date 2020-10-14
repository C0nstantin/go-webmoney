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
	invoice := webmoney.Invoice{
		OrderId:      "111",
		CustomerWmid: "128756507061",
		StorePurse:   "R697541065597",
		Amount:       "3.2",         //AMOUNT
		Desc:         "DESCRIPTION", //DESCRIPTION FOR INVOICE
		Address:      "SHOP_ADDRESS",
		Period:       "1",
		Expiration:   "1",
		OnlyAuth:     "0",
	}

	result, err := wmClient.SendInvoice(invoice)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
}
