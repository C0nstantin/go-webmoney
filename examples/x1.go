//Webmoney XML Interfaces
//Interface X1. Sending invoice from merchant to customer.

package main

import "github.com/C0nstantin/webmoney-go/pkg/webmoney"

func main() {
	wmClient := webmoney.WmClient{
		Wmid: "YOUR_WMID",
		Key:  "SECRET_KEY",
		Pass: "SECRET_PASSWORD",
	}

	invoice := webmoney.Invoice{
		OrderId:      "ORDER ID IN YOUR SISTEM",
		CustomerWmid: "CUSTOMER WMID",
		StorePurse:   "MERCHANT PURSE",
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

	fmt(result.Id)
}
