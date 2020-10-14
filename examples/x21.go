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
	got := webmoney.X21Request{
		LmiPayeePurse:       "Z303339773989",
		LmiDayLimit:         "1",
		LmiWeekLimit:        "1",
		LmiMonthLimit:       "1",
		LmiClientnumber:     "128756507061",
		LmiClientnubmerType: "1",
		LmiSmsType:          "1",
		Lang:                "ru-RU",
	}

	result, err := wmClient.TrustRequest(got)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
	fmt.Println("Insert code for confirmation")
	var code string
	fmt.Scan(&code)
	got1 := webmoney.X212Request{
		LmiPurseId:          result.TrustX21.PurseId,
		LmiClientnumberCode: code,
	}

	result1, err := wmClient.TrustConfirm(got1)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result1)
}
