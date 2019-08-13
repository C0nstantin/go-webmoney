// Copyright 2015 Constantin Karataev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package provides Webmoney xml interfaces
// References
//    https://wiki.wmtransfer.com/projects/webmoney/wiki/XML-interfaces
//    https://wiki.webmoney.ru/projects/webmoney/wiki/XML-%D0%B8%D0%BD%D1%82%D0%B5%D1%80%D1%84%D0%B5%D0%B9%D1%81%D1%8B

package webmoney

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/C0nstantin/go-webmoney/wmsigner"
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// cononical name for structs
const (
	X1   string = `Invoice`
	X2   string = `Trans`
	X3   string = `Operations`
	X4   string = `OutInvoices`
	X5   string = `FinishProtect`
	X6   string = `SendMsg`
	X7   string = `ClassicAuth`
	X8   string = `FindWMPurseNew`
	X9   string = `Purses`
	X10  string = `InInvoices`
	X11  string = `GetWMPassport`
	X13  string = `RejectProtect`
	X14  string = `TransMoneyback`
	X15  string = `TrustList`
	X152 string = `TrustList2`
	X153 string = `TrustSave2`
	X16  string = `CreatePurse`
	X23  string = `InvoiceRefusal`
)

// base scruct for all request
type requestBase struct {
	XMLName xml.Name    `xml:"w3s.request"`
	Reqn    string      `xml:"reqn"`
	Wmid    string      `xml:"wmid"`
	Sign    string      `xml:"sign"`
	Request interface{} `xml:",>"`
}

//base struct for all response
type responseBase struct {
	XMLName  xml.Name    `xml:"w3s.response"`
	Reqn     string      `xml:"reqn"`
	Retval   int64       `xml:"retval"`
	Retdesc  string      `xml:"retdesc"`
	Ser      string      `xml:"ser"`
	CWmid    string      `xml:"cwmid"`
	Wmid     string      `xml:"wmid"`
	Response interface{} `xml:",any"`
}

//Struct for initicalize Webmoney client and save response
type WmClient struct {
	Wmid      string
	Key       string
	Pass      string
	Cert      string
	Sign      string
	Reqn      string
	X         string
	ResultStr string
	Request   interface{}
	Result    responseBase
}

// Function return true if current settings indicate
// that request is signed classic key
func (w *WmClient) IsClassic() bool {
	if w.Key != "" && w.Pass != "" {
		return true
	} else {
		return false
	}
}

// Function return true if current settings indicate
// that reuests is signed light keeper
func (w *WmClient) IsLight() bool {
	if w.Key != "" && w.Cert != "" {
		return true
	} else {
		return false
	}
}

// Function check settings for connetion and sign not set
// before start use you must set Wmid, Key and Pass for Keeper Classic(WinPro)
// or wmid, key and cert for Keepr Light(WebPro)
func (w *WmClient) noInit() bool {
	if w.Wmid == "" || w.Key == "" || w.Pass == "" {
		return true
	} else {
		return false
	}
}

// Function create xml request and send in interface url
func (w *WmClient) SendRequest() (string, error) {
	v := &requestBase{
		Wmid: w.Wmid,
		Reqn: w.Reqn,
	}

	if w.IsClassic() && w.Sign != "" {
		s := wmsigner.NewSigner(w.Wmid, w.Pass, w.Key)
		if result, err := s.Sign(w.Sign); err != nil {
			log.Fatal(err)
		} else {
			v.Sign = result
		}
	} else {
		v.Sign = ""
	}

	var url string

	if w.IsClassic() {
		url = "https://w3s.webmoney.ru/asp/XML" + w.X + ".asp"
	} else {
		if w.X != X8 {
			url = "https://w3s.wmtransfer.com/asp/XML" + w.X + "Cert.asp"
		} else {
			url = "https://w3s.wmtransfer.com/asp/XML" + w.X + "XMLFindWMPurseCertNew.asp"
		}
	}
	v.Request = w.Request

	output, err := xml.MarshalIndent(v, "  ", "    ")
	body := "<?xml version=\"1.0\" encoding=\"utf-8\"?> \n" + string(output)

	result, err := w.sendRequest(url, body)
	return result, err
}

// private
// Functrion send requst to server and return response how string
func (w *WmClient) sendRequest(url string, body string) (string, error) {
	var tr *http.Transport
	// load root ca
	r := strings.NewReader(ROOT_CA)
	caCert, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	if w.IsClassic() {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
			DisableCompression: true,
		}
	} else {
		//load client ca
		cert, err := tls.LoadX509KeyPair(w.Cert, w.Key)
		if err != nil {
			log.Fatal(err)
		}
		tlsConfig := &tls.Config{
			Certificates:  []tls.Certificate{cert},
			RootCAs:       caCertPool,
			Renegotiation: tls.RenegotiateFreelyAsClient, //RenegotiateOnceAsClient,
		}
		tlsConfig.BuildNameToCertificate()

		tr = &http.Transport{
			TLSClientConfig:    tlsConfig,
			DisableCompression: true,
		}
	}
	fmt.Println(body)
	client := &http.Client{Transport: tr}
	resp, err := client.Post(url, "text/xml", strings.NewReader(body))

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	w.ResultStr = string(result)
	fmt.Println(w.ResultStr)
	if err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

//Function prse response and return structure response
func (w *WmClient) ParseResponse(resp interface{}) error {
	v := responseBase{
		Response: resp,
	}

	r := bytes.NewReader([]byte(w.ResultStr))
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charset.NewReader
	err := dec.Decode(&v)
	w.Result = v
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (w *WmClient) getResult(result interface{}) error {
	if _, err := w.SendRequest(); err != nil {
		return err
	}
	if err := w.ParseResponse(result); err != nil {
		return err
	}
	if (w.Result.Retval != 0) && (w.Result.Retval != 1 && w.X == X8) {
		err := errors.New(strconv.FormatInt(w.Result.Retval, 10) + "   " + w.Result.Retdesc)
		return err
	}

	return nil

}
