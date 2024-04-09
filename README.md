go-webmoney 

# Requirements

* Go1.7

ABOUT
------
XML interfaces are basically text message interfaces that send data via HTTPS protocol (http over 128-bit SSL) to special certification web servers of the system. Data is sent in the XML format with help of special module for authentication of WebMoney Keeper key files or standard certificates (WM Keeper WebPro certificates).
[Detail](http://wiki.wmtransfer.com/projects/webmoney/wiki/XML-interfaces)

# Install
```go 
package main

import "github.com/C0nstantin/go-webmoney"
//...
```

How Use
-------

# WinPro Client (Classic)
```go
var CnfClassic = struct {
  Wmid string
  Key  string
  Pass string
}{
  `000000000000`, //service wmid
  `---key in base64 string---`,
  `password for key`,
}

wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Key,
		Pass: CnfClassic.Pass,
	}

```
# WebPro Client (Light)
```go
var CnfLight = struct{
  Wmid string
  Cert string
  Key  string    
}{
   `000000000000`, //service wmid
  "/path/to/cert/file.pem",
  `/path/to/key/file.key`,
   
}

wmCl := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Cert:  CnfLight.Cert,
		Key: CnfLight.Key,
}

```
#WM Signer

```go
func main(){
	str := "source string"
	signer := wmsigner.NewSigner("wmid","password","key_string")
	sign_str,err  := signer.Sign(str)
	if(err != nil) {
		// do something if error})
	}

	//use sign_str
```


# Webmoney XML Interfaces

##X1 Sending Invoice from merchant to customer

[Detail Desctiption Interface](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X1)
```go
invoice := webmoney.Invoice{
		OrderId:      "10",
		CustomerWmid: "000000000000", //cutomer wmid
		StorePurse:   "R123456789098", //store purse
		Amount:       "3.2",
		Desc:         "description invoice",
		Address:      "Address for delivery",
		Period:       "1",
		Expiration:   "1",
		OnlyAuth:     "0",
	}
	result, err := wmCl.SendInvoice(invoice)
/*
result implements type:
  type InvoiceResponse struct {
    Id           string `xml:"id,attr"`
    Ts           string `xml:"ts,attr"`
    OrderId      string `xml:"orderid"`
    CustomerWmid string `xml:"customerwmid"`
    StorePurse   string `xml:"storepurse"`
    Amount       string `xml:"amount"`
    Desc         string `xml:"desc"`
    Address      string `xml:"address"`
    Period       string `xml:"period"`
    Expiration   string `xml:"expiration"`
    State        string `xml:"state"`
    DateCrt      string `xml:"datecrt"`
    DateUpd      string `xml:"dateupd"`
    WmTranId     string `xml:"wmtranid"`
  }

*/

```
## X2 Transferring funds from one purse to another.
This interface is available for registered members only and can be used for making 
transfers from purses of any WM Keeper, including Budget Automates.
**The option can be enabled at the Web Merchant Interface service settings page in the "additional parameters" section. 
Enabling this option you undertake to use the "trans" nlyauth = 1" parameter.
The Keeper that signs requests should have Personal passport. "**
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X2)
```go
	transaction := webmoney.Transaction{
		TranId:    "1001",
		PurseDest: "R000000000000",
		PurseSrc:  "R111111111111",
		Amount:    "3.2",
		Period:    "1",
		Desc:      "Теst кирилица",
		PCode:     "",
		WmInvid:   "0",
		OnlyAuth:  "1",
	}

	result, err := wmClientClassic.CreateTransaction(transaction)
/*
result implements:
  
  type Operation struct {
    XMLName   xml.Name `xml:operation`
    Id        string   `xml:"id,attr"`
    Ts        string   `xml:"ts,attr"`
    TranId    string   `xml:"tranid"`
    PurseSrc  string   `xml:"pursesrc"`
    PurseDest string   `xml:"pursedest"`
    Amount    string   `xml:"amount"`
    Commis    string   `xml:"comiss"`
    Opertype  string   `xml:"opertype"`
    Period    string   `xml:"period"`
    WmInvid   string   `xml:"wminvid"`
    Desc      string   `xml:"desc"`
    DateCrt   string   `xml:"datecrt"`
    DateUpd   string   `xml:"dateupd"`
    CorrWm    string   `xml:"corrwm"`
    Rest      string   `xml:"rest"`
    TimeLock  bool     `xml:timelock`
  }
*/
```
## x3 Receiving the History of Transactions; checking Transaction status
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X3)
```go
	opts := webmoney.GetOpers{
		TranId:     "10",// optional
		Purse:      "R11111111111",
		DateStart:  "20160726 00:00:00",
		DateFinish: "20160721 00:00:00",
	}
  result, err := wmCl.GetOperations(opts)
  /*
  result implements 
  type Operations struct {
	XMLName       xml.Name    `xml:"operations"`
	OperationList []Operation `xml:"operation"` // x1 collection
}
  */
```
## x4 Receiving the history of issued invoices. Verifying whether invoices were paid
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X4)

```go
invoices := webmoney.OutInvoices{
		Purse:      "R111111111111",
		DateStart:  "20160622 00:00:00",
		DateFinish: "20160927 00:00:00",
	}
	result, err := wmCl.GetOutInvoices(invoices)
/*
result implements

type OutInvoicesResp struct {
	XMLName     xml.Name          `xml:"outinvoices"`
	InvoiceList []InvoiceResponse `xml:"outinvoice"` //x1
}
*/
```
## x5 Completing a code-protected transaction. Entering a protection code
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X5)
```go
//create transactionf with pcode

	wmCl := webmoney.WmClient{
		Wmid: CnfClassic.Wmid,
		Key:  CnfClassic.Wmid_key,
		Pass: CnfClassic.Wmid_pass,
	}
	pcode := "testtest"
	tt := webmoney.Transaction{
		TranId:    "1999",
		PurseDest: "R000000000000",
		PurseSrc:  "R222222222222",
		Amount:    "3.2",
		Period:    "1",
		Desc:      "Desc for test",
		WmInvid:   "0",
		OnlyAuth:  "1",
		PCode:     pcode,
	}

	result, err := wmCl.CreateTransaction(tt)

	if err != nil {
		log.Fatal(err)
	}
 //registration second client 
	wmClight := webmoney.WmClient{
		Wmid: CnfLight.Wmid,
		Cert: CnfLight.Cert,
		Key:  CnfLight.Key,
	}

	fp := webmoney.FinishProtect{
		WmTranId: result.Ts,
		PCode:    pcode,
	}
  //finishing transaction
	resultFinish, err := wmClight.DoFinishProtect(fp)
  /* resultFinish implements type Operation */
```
## x6 Sending message to random WM-identifier via internal mail
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X6)
```go

	wmMsg := webmoney.SendMsg{
		ReceiverWmid: "00000000000",
		MsgSubj:      "testClassic",
		MsgText:      "Тест сообщения через Классик кирилица",
	}

	result, err := wmCl.SendMessage(wmMsg)
  /*result implements
    type SendMsgResponse struct {
      XMLName      xml.Name `xml:"message"`
      MessageId    int      `xml:"id,attr"`
      ReceiverWmid string   `xml:"receiverwmid"`
      MsgSubj      string   `xml:"msgsubj"`
      MsgText      string   `xml:"msgtext"`
      DateCrt      string   `xml:"datecrt"`
    }

  */

```

## x7 Verifying client’s handwritten signature – owner of WM Keeper WinPro
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X7)
```go
	signature := webmoney.TestSignRequest{
		Wmid: "128756507061",
		Plan: "plantext",
		Sign: "signtext",
	}
	result, err := wmCl.TestSign(signature)
  /*
  result implements
    type TestSignResponse struct {
  	  XMLName xml.Name `xml:"testsign"`
	    Res     string   `xml:"res"`
    }
  */
```

## x8 Retrieving information about purse ownership. Searching for system user by his/her identifier or purse
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X8)
```go

	testWm := webmoney.TestWmPurse{
		Wmid: `000000000000`,
	}

	testPurse := webmoney.TestWmPurse{
		Purse: `Z0000000000000`,
	}
	testWmidPurse := webmoney.TestWmPurse{
		Wmid:  `0000000000000`,
		Purse: `K000000000000`,
	}

	result, err := wmCl.FindWmidPurse(testWm)
	
  //OR
  
  result, err := wmCl.FindWmidPurse(testPurse)
	
  //OR
  
  result, err := wmCl.FindWmidPurse(testWmidPurse)
	/*result implements
  type TestWmPurseResponse struct {
	  XMLName xml.Name      `xml:"testwmpurse"`
	  Wmid    Wmid          `xml:"wmid"`
	  Purse   ReturnedPurse `xml:"purse"`
  }
  type ReturnedPurse struct {
	  XMLName              xml.Name `xml:"purse"`
	  Value                string   `xml:",chardata"`
	  MerchantActiveMode   string   `xml:"merchant_active_mode,attr"`
	  MerchantAllowCashier string   `xml:"merchant_allow_cashier,attr"`
  }
  type Wmid struct {
	  XMLName           xml.Name `xml:"wmid"`
	  Value             string   `xml:",chardata"`
	  Available         string   `xml:"available,attr"`
	  Themselfcorrstate string   `xml:"themselfcorrstate,attr"`
	  Newattst          string   `xml:"newattst,attr"`
  }
  */
```
## x9 Retrieving information about purse balance
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X9)
```go
  	p := webmoney.Purses{
		Wmid: CnfClassic.Wmid,
	}
	result, err := wmCl.GetPurses4Wmid(p)
  /*result implements
    type PursesResp struct {
	    XMLName   xml.Name `xml:"purses"`
	    Cnt       string   `xml:"cnt,attr"`
	    PurseList []RPurse `xml:"purses"`
    }

    type RPurse struct {
	    XMLName     xml.Name `xml:"purse"`
	    PurseName   string   `xml:"pursename"`
	    Amount      string   `xml:"amount"`
	    Desc        string   `xml:"desc"`
	    Outsideopen string   `xml:"outsideopen"`
	    LastIntr    string   `xml:"lastintr"`
	    LastOuttr   string   `xml:"lastouttr"`
    }
  */
```

## x10 Retrieving list of invoices for payment
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X10)
```go
	i := webmoney.InInvoices {
		Wmid:       CnfLight.Wmid,
		WmInvid:    "0",
		DateStart:  "20140101 00:00:00",
		DateFinish: "20160930 00:00:00",
	}
	result, err := wmCl.GetInInvoices(i)
  /*result implements
    type InInvoicesResponse struct {
	    XMLName     xml.Name          `xml:"ininvoices"`
	    InvoiceList []InvoiceResponse `xml:ininvoices`
    }

  */
```

## x11 Retrieving information from client’s passport by WM-identifier
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X11)
```go
	result, err := webmoney.GetInfoWmid(`128756507061`)
 ```

## x13 Recalling incomplete protected transaction
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X13)
```go
  rp := RejectProtect {
    WmTranId: "111212",
  }
  result, err := WmCl.RejectProtect(rp)
  /*result implements type Operation*/
```

## x14 Fee-free refund
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X14)
```go
   tr := webmoney.Trans {
  	InWmTranId:"234",
	  Amount: "242",
	  MoneyBackPhone:"925000000",
  }
  result,err := wmCl.MoneyBack(tr)
  /*result implements type Operation*/

```

## x15 Viewing and changing settings of “by trust” management
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X15)
```go

//get trastlist
tl := webmoney.Trustlist {
    Wmid:"000000000000",
  }
result,err := wmCl.GetTrustList(tl)
/*
result implemenst

type TrustListResponse struct {
	XMLName   xml.Name `xml:"trustlist"`
	Cnt       string   `xml:"cnt,attr"`
	TrustList []Trust  `xml:"trastlist"`
}

type Trust struct {
	XMLName     xml.Name `xml:"trast"`
	Id          string   `xml:"id,attr"`
	Inv         string   `xml:"inv,attr"`
	Trast       string   `xml:"trans,attr"`
	PurseAttr   string   `xml:"purse,attr"`
	TransHist   string   `xml:"transhist"`
	Master      string   `xml:"master"`
	Purse       string   `xml:"purse"`
	DayLimit    string   `xml:"daylimit"`
	DLimit      string   `xml:"dlimit"`
	WLimit      string   `xml:"wlimit"`
	MLimit      string   `xml:"mlimit"`
	DSum        string   `xml:"dsum"`
	WSum        string   `xml:"wsum"`
	MSum        string   `xml:"msum"`
	LastSumDate string   `xml:lastsumdate`
}

*/

st := webmoney.TrastSave{
    Inv: "",
    Trans:"",
    PurseAttr:"",
    TransHist:"",
    MasterWmid:"",
    SlaveWmid:"",
    Purse:"",
    Limit:"",
    DayLimit:"",
    WeekLimit:"",
    MonthLimit:"",
  }
result,err := wmCl.SetTrast(st)
/* result implements Trust */

```

## x16 Creating a purse
[Detail Information](http://wiki.wmtransfer.com/projects/webmoney/wiki/Interface_X16)
```go
p := webmoney.Purse{
    Wmid:"111111111110",
    PurseType:"Z",
    Desc:"test",
  }
result,error := wmCl.CreatePurese(p)
/*
result implements 

type PurseResponse struct {
	XMLName   xml.Name `xml:"purse"`
	Id        string   `xml:"id,attr"`
	PurseName string   `xml:"pursename"`
	Amount    string   `xml:"amount"`
	Desc      string   `xml:"desc"`
}


*/

```


##x17-x23
see examples file

## Custom domains
If you need change urls for api requests please use this environment variables 
```go
//for merchant requests
os.Setenv("MERCHANT_DOMAIN", "merchant.web.money") //merchant.web.money is a default value
//if you use classic (winpro) auth
os.Setenv("CLASSIC_DOMAIN", "w3s.webmoney.com") // w3s.webmoney.com is a default value
//if you use light (webpro) auth
os.Setenv("LIGHT_DOMAIN", "w3s.webmoney.com") // w3s.webmoney.com is a default value
// for passport request 
os.Setenv("API_PASSPORT_DOMAIN", apipassort.web.money)
```
