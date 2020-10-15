// Copyright 2015 Constantin Karataev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package provides Webmoney xml interfaces
// References
//    https://wiki.wmtransfer.com/projects/webmoney/wiki/XML-interfaces
//    https://wiki.webmoney.ru/projects/webmoney/wiki/XML-%D0%B8%D0%BD%D1%82%D0%B5%D1%80%D1%84%D0%B5%D0%B9%D1%81%D1%8B

package webmoney

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"strings"

	_ "github.com/paulrosania/go-charset/data"
)

type WmClient struct {
	Wmid      string
	Key       string
	Pass      string
	Cert      string
	SecretKey string
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

// private
// Functrion send requst to server and return response how string
func (w *WmClient) sendRequest(url string, body string) (string, error) {
	tr, err := w.getTransport()
	if err != nil {
		return "", err
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Post(url, "text/xml", strings.NewReader(body))
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	} else {
		return string(result), nil
	}
}

func (w *WmClient) getTransport() (*http.Transport, error) {
	var tr *http.Transport
	// load root ca
	r := strings.NewReader(ROOT_CA)
	caCert, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	if w.IsClassic() {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{
				//RootCAs: caCertPool,
				Renegotiation: tls.RenegotiateFreelyAsClient,
			},
			DisableCompression: true,
		}
	} else {

		cert, err := tls.LoadX509KeyPair(w.Cert, w.Key)
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			Certificates:  []tls.Certificate{cert},
			//RootCAs:       caCertPool,
			Renegotiation: tls.RenegotiateFreelyAsClient,
		}
		tlsConfig.BuildNameToCertificate()

		tr = &http.Transport{
			TLSClientConfig:    tlsConfig,
			DisableCompression: true,
		}
	}
	return tr, nil
}
