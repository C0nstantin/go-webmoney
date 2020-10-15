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

	test := webmoney.InInvoices{

		Wmid:       "128756507061",
		WmInvid:    "0",
		DateStart:  "20201012",
		DateFinish: "20201013",
	}

	result, err := wmClient.GetInInvoices(test)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
}
