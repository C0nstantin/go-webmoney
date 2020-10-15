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

	msg := webmoney.SendMsg{
		ReceiverWmid: "128756507061",
		MsgText: `Тестовое сообщение многострочное
сообщение`,
		MsgSubj:  "Тестовое сообещие",
		OnlyAuth: 0,
	}

	result, err := wmClient.SendMessage(msg)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
}
