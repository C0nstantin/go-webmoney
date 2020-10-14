package main

import (
	"fmt"
	"github.com/C0nstantin/go-webmoney"
	"log"
)

func main() {

	result, err := webmoney.GetInfoWmid("128756507061")

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%#v", result)
}
