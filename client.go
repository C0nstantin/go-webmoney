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
	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"webmoney/wmsigner"
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
// Root ca for webmoney requests
const ROOT_CA = `
-----BEGIN CERTIFICATE-----
MIIFsTCCA5mgAwIBAgIQA7dHzSZ7uJdBxFycIWn+WjANBgkqhkiG9w0BAQUFADBr
MSswKQYDVQQLEyJXTSBUcmFuc2ZlciBDZXJ0aWZpY2F0aW9uIFNlcnZpY2VzMRgw
FgYDVQQKEw9XTSBUcmFuc2ZlciBMdGQxIjAgBgNVBAMTGVdlYk1vbmV5IFRyYW5z
ZmVyIFJvb3QgQ0EwHhcNMTAwMzEwMTczNDU2WhcNMzUwMzEwMTc0NDUxWjBrMSsw
KQYDVQQLEyJXTSBUcmFuc2ZlciBDZXJ0aWZpY2F0aW9uIFNlcnZpY2VzMRgwFgYD
VQQKEw9XTSBUcmFuc2ZlciBMdGQxIjAgBgNVBAMTGVdlYk1vbmV5IFRyYW5zZmVy
IFJvb3QgQ0EwggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDFLJXtzEkZ
xLj1HIj9EhGvajFJ7RCHzF+MK2ZrAgxmmOafiFP6QD/aVjIexBqRb8SVy38xH+wt
hqkZgLMOGn8uDNpFieEMoX3ZRfqtCcD76KDySTOX1QUwHAzBfGuhe2ZQULUIjxdP
Ra4NDyvmXh4pE/s1+/7dGbUZs/JpYYaD2xxAt5PDTjylsKOk4FMb5kv6jzORkXku
5UKFGUXEXbbf1xzgYHMIzoeJGn+iPgVFYAvkkQyvxEaVj0lNE+q/ua761krgCo47
BiH1zMFzkv4uNHEZfe/lyHaozzbsu6yaK3EdrURSLuWrlxKy9yo3xDe9TPkzkhPe
JPbV7YgvUUtWSeAJpksBU8GCALEhSgXOfHckuJFj9QB3YecHBvjdSiAUuntwM/iH
vtSOXEUHxqW75E2Gq/2L4vBcxArXVdbUrVQDF3klzYu17OFgfe1hHHMHzgr4HBML
ZiRCcvNLqghBCVxu1DM15YDfw+wnNV/5dUPx60tiocmCZpJKTwVl8gc85QCPyREu
jey8F0kgdgssQosPWTTWDg7X4Ifw20VkplHZDr29K5HdwLe56TvOI/4H24XJdqpA
xoLBx9PL6ZXxH52wU0bSluL8/joXGzavFrhsXH7jJocH6tsFVzBZrmnVswbUMHDN
L3xSnr5fAAXXZa7UwHd3pq/fsdG7s9PByQIDAQABo1EwTzALBgNVHQ8EBAMCAYYw
DwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUsTCnSwOZT4Q2HBN9V/TrafuIG8Mw
EAYJKwYBBAGCNxUBBAMCAQAwDQYJKoZIhvcNAQEFBQADggIBAAy5jHDFpVWtF209
N30I+LHsiqMaLUmYDeV6sUBJqmWAZav7pWnigiMLuJd9yRa/ow6yKlKPRi3sbKaB
wsAQ+xnz811nLFBBdS4PkwlHu1B7P4B2YbcqmF6k1QieJBZxOn7wledtnoBAkZ4d
6HEW1OM5cvCoyj8YAdJTZIBzn61aNt/viPvypIUQf6Ps6Q2daNEAj7DoxIY8crnO
aSIGdGmlT/y/edSqWv9Am5e9KXkJhQWMnGXh43wJYyHTetxVWPS43bW7gIUADYyc
KSH3isrBN5xQOFXMfL+lVHHSs7ap23DOo7xIDenm5PWz+QdDDFz3zLVeRovnkIdk
a/Wgk3f6rFfKB0y5POJ+BJvkorIYNZiN3dnmc6cDP840BUMv3BUrOe8iSy5lRr8m
R+daktbZfO8E/rAb3zEdN+KG/CNJfAnQvp6DT4LqY/J9pG+VusH5GpUwuXr7UqLw
End1LRp7qm28Cic7fegUnnUpkuF4ZFq8pWq8w59sOWlRuKBuWX46OghMrjgD0AN1
hlA2/d5ULImX70Q2te3xiS1vrQhu77mkb/jA4/9+YfeT7VMpbnC3OoHiZ2bjudKn
thlOs+AuUvzB4Tqo62VSF5+r0sYI593S+STmaZBAzsoaoEB7qxqKbEKCvXb9BlXk
L76xIOEkbSIdPIkGXM4aMo4mTVz7
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIENjCCAx6gAwIBAgIBATANBgkqhkiG9w0BAQUFADBvMQswCQYDVQQGEwJTRTEU
MBIGA1UEChMLQWRkVHJ1c3QgQUIxJjAkBgNVBAsTHUFkZFRydXN0IEV4dGVybmFs
IFRUUCBOZXR3b3JrMSIwIAYDVQQDExlBZGRUcnVzdCBFeHRlcm5hbCBDQSBSb290
MB4XDTAwMDUzMDEwNDgzOFoXDTIwMDUzMDEwNDgzOFowbzELMAkGA1UEBhMCU0Ux
FDASBgNVBAoTC0FkZFRydXN0IEFCMSYwJAYDVQQLEx1BZGRUcnVzdCBFeHRlcm5h
bCBUVFAgTmV0d29yazEiMCAGA1UEAxMZQWRkVHJ1c3QgRXh0ZXJuYWwgQ0EgUm9v
dDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALf3GjPm8gAELTngTlvt
H7xsD821+iO2zt6bETOXpClMfZOfvUq8k+0DGuOPz+VtUFrWlymUWoCwSXrbLpX9
uMq/NzgtHj6RQa1wVsfwTz/oMp50ysiQVOnGXw94nZpAPA6sYapeFI+eh6FqUNzX
mk6vBbOmcZSccbNQYArHE504B4YCqOmoaSYYkKtMsE8jqzpPhNjfzp/haW+710LX
a0Tkx63ubUFfclpxCDezeWWkWaCUN/cALw3CknLa0Dhy2xSoRcRdKn23tNbE7qzN
E0S3ySvdQwAl+mG5aWpYIxG3pzOPVnVZ9c0p10a3CitlttNCbxWyuHv77+ldU9U0
WicCAwEAAaOB3DCB2TAdBgNVHQ4EFgQUrb2YejS0Jvf6xCZU7wO94CTLVBowCwYD
VR0PBAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wgZkGA1UdIwSBkTCBjoAUrb2YejS0
Jvf6xCZU7wO94CTLVBqhc6RxMG8xCzAJBgNVBAYTAlNFMRQwEgYDVQQKEwtBZGRU
cnVzdCBBQjEmMCQGA1UECxMdQWRkVHJ1c3QgRXh0ZXJuYWwgVFRQIE5ldHdvcmsx
IjAgBgNVBAMTGUFkZFRydXN0IEV4dGVybmFsIENBIFJvb3SCAQEwDQYJKoZIhvcN
AQEFBQADggEBALCb4IUlwtYj4g+WBpKdQZic2YR5gdkeWxQHIzZlj7DYd7usQWxH
YINRsPkyPef89iYTx4AWpb9a/IfPeHmJIZriTAcKhjW88t5RxNKWt9x+Tu5w/Rw5
6wwCURQtjr0W4MHfRnXnJK3s9EK0hZNwEGe6nQY1ShjTK3rMUUKhemPR5ruhxSvC
Nr4TDea9Y355e6cJDUCrat2PisP29owaQgVR1EX1n6diIWgVIEM8med8vSTYqZEX
c4g/VhsxOBi0cQ+azcgOno4uG+GMmIPLHzHxREzGBHNJdmAPx/i9F4BrLunMTA5a
mnkPIAou1Z5jJh5VkpTYghdae9C8x49OhgQ=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIICPDCCAaUCEHC65B0Q2Sk0tjjKewPMur8wDQYJKoZIhvcNAQECBQAwXzELMAkG
A1UEBhMCVVMxFzAVBgNVBAoTDlZlcmlTaWduLCBJbmMuMTcwNQYDVQQLEy5DbGFz
cyAzIFB1YmxpYyBQcmltYXJ5IENlcnRpZmljYXRpb24gQXV0aG9yaXR5MB4XDTk2
MDEyOTAwMDAwMFoXDTI4MDgwMTIzNTk1OVowXzELMAkGA1UEBhMCVVMxFzAVBgNV
BAoTDlZlcmlTaWduLCBJbmMuMTcwNQYDVQQLEy5DbGFzcyAzIFB1YmxpYyBQcmlt
YXJ5IENlcnRpZmljYXRpb24gQXV0aG9yaXR5MIGfMA0GCSqGSIb3DQEBAQUAA4GN
ADCBiQKBgQDJXFme8huKARS0EN8EQNvjV69qRUCPhAwL0TPZ2RHP7gJYHyX3KqhE
BarsAx94f56TuZoAqiN91qyFomNFx3InzPRMxnVx0jnvT0Lwdd8KkMaOIG+YD/is
I19wKTakyYbnsZogy1Olhec9vn2a/iRFM9x2Fe0PonFkTGUugWhFpwIDAQABMA0G
CSqGSIb3DQEBAgUAA4GBALtMEivPLCYATxQT3ab7/AoRhIzzKBxnki98tsX63/Do
lbwdj2wsqFHMc9ikwFPwTtYmwHYBV4GSXiHx0bH/59AhWM1pF+NEHJwZRDmJXNyc
AA9WjQKZ7aKQRUzkuxCkPfAyAw7xzvjoyVGM5mKf5p/AfbdynMk2OmufTqj/ZA1k
-----END CERTIFICATE-----`
