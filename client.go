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
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/paulrosania/go-charset/data"
)

// Default timeout for WebMoney API requests
var defaultHTTPTimeout = 30 * time.Second

type WmClient struct {
	Wmid      string
	Key       string
	Pass      string
	Cert      string
	SecretKey string

	mu         sync.Mutex
	httpClient *http.Client
}

func maskSecret(s string) string {
	if len(s) == 0 {
		return ""
	}
	return "****"
}

// String masks sensitive fields to prevent accidental secret leakage in logs.
func (w WmClient) String() string {
	return fmt.Sprintf("{Wmid:%s Key:%s Pass:%s Cert:%s SecretKey:%s}",
		w.Wmid, maskSecret(w.Key), maskSecret(w.Pass), w.Cert, maskSecret(w.SecretKey))
}

// GoString masks sensitive fields for fmt %#v formatting.
func (w WmClient) GoString() string {
	return fmt.Sprintf("webmoney.WmClient{Wmid:%q, Key:%q, Pass:%q, Cert:%q, SecretKey:%q}",
		w.Wmid, maskSecret(w.Key), maskSecret(w.Pass), w.Cert, maskSecret(w.SecretKey))
}

// MarshalJSON masks sensitive fields to prevent secret leakage through JSON serialization.
func (w WmClient) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Wmid      string `json:"Wmid"`
		Key       string `json:"Key"`
		Pass      string `json:"Pass"`
		Cert      string `json:"Cert"`
		SecretKey string `json:"SecretKey"`
	}{
		Wmid:      w.Wmid,
		Key:       maskSecret(w.Key),
		Pass:      maskSecret(w.Pass),
		Cert:      w.Cert,
		SecretKey: maskSecret(w.SecretKey),
	})
}

// IsClassic function return true if current settings indicate
// that request is signed classic key
func (w *WmClient) IsClassic() bool {
	if w.Key != "" && w.Pass != "" {
		return true
	} else {
		return false
	}
}

// IsLight function return true if current settings indicate
// that requests is signed light keeper
func (w *WmClient) IsLight() bool {
	if w.Key != "" && w.Cert != "" {
		return true
	} else {
		return false
	}
}

// noInit check settings for connetion and sign not set
// before start use you must set Wmid, Key and Pass for Keeper Classic(WinPro)
// or wmid, key and cert for Keepr Light(WebPro)
func (w *WmClient) noInit() bool {
	if w.Wmid == "" || w.Key == "" || w.Pass == "" {
		return true
	} else {
		return false
	}
}

// getHTTPClient returns a thread-safe cached http.Client instance.
// This prevents connection leaks by ensuring the transport is created only once.
func (w *WmClient) getHTTPClient() (*http.Client, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.httpClient != nil {
		return w.httpClient, nil
	}

	if os.Getenv("USE_W3S_CERT") != "" {
		tr, err := w.getTransport()
		if err != nil {
			return nil, err
		}
		w.httpClient = &http.Client{
			Transport: tr,
			Timeout:   defaultHTTPTimeout,
		}
	} else {
		w.httpClient = &http.Client{
			Timeout: defaultHTTPTimeout,
		}
	}
	
	return w.httpClient, nil
}

// Function send requst to server and return response how string
func (w *WmClient) sendRequest(url string, body string) (string, error) {
	client, err := w.getHTTPClient()
	if err != nil {
		return "", err
	}

	resp, err := client.Post(url, "text/xml", strings.NewReader(body))
	if err != nil {
		return "", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println(err)
		}
	}(resp.Body)
	result, err := io.ReadAll(resp.Body)
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
	caCert, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	if w.IsClassic() {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{
				//RootCAs: caCertPool,
			},
			DisableCompression: true,
		}
	} else {

		cert, err := tls.LoadX509KeyPair(w.Cert, w.Key)
		if err != nil {
			return nil, err
		}
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			//RootCAs:       caCertPool,
			Renegotiation: tls.RenegotiateFreelyAsClient, // WebMoney light certificates still need this
		}

		tr = &http.Transport{
			TLSClientConfig:    tlsConfig,
			DisableCompression: true,
		}
	}
	return tr, nil
}
