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

	transaction := webmoney.Transaction{
		PurseSrc:  "Z30333977000",
		PurseDest: "Z214605808000",
		Amount:    "0.01",
		OnlyAuth:  "0",
		TranId:    "4",
		WmInvid:   "0",
		Period:    "2",
		PCode:     "11111",
		Desc:      "test",
	}

	result, err := wmClient.CreateTransaction(transaction)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(result.Id)
}
