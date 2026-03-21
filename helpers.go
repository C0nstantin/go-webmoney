// Copyright 2015 Constantin Karataev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// This file contains some helper function to form data for requests

package webmoney

import (
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

// maxResponseSize limits the maximum WebMoney API response size to 10 MB.
// This prevents memory exhaustion from a malicious or misconfigured server.
const maxResponseSize = 10 * 1024 * 1024

// validDomainRegex matches valid domain names (letters, digits, hyphens, dots).
var validDomainRegex = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?)+$`)

// isValidDomain checks that a string looks like a valid domain name.
// Rejects paths, query strings, ports, and other URL components.
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}
	return validDomainRegex.MatchString(domain)
}

var sensitiveTags = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(<secret_key>)(.*?)(</secret_key>)`),
	regexp.MustCompile(`(?i)(<pcode>)(.*?)(</pcode>)`),
	regexp.MustCompile(`(?i)(<sign>)(.*?)(</sign>)`),
	regexp.MustCompile(`(?i)(<md5>)(.*?)(</md5>)`),
	regexp.MustCompile(`(?i)(<sha256>)(.*?)(</sha256>)`),
}

// MaskXMLSecrets replaces content inside sensitive XML tags with asterisks.
func MaskXMLSecrets(xmlStr string) string {
	result := xmlStr
	for _, re := range sensitiveTags {
		result = re.ReplaceAllString(result, `$1****$3`)
	}
	return result
}

var (
	moscowLoc *time.Location
	utcLoc    *time.Location
	locOnce   sync.Once
)

func initLocations() {
	var err error
	moscowLoc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		// Fallback for systems without tzdata: Moscow is UTC+3
		moscowLoc = time.FixedZone("MSK", 3*60*60)
	}
	utcLoc, err = time.LoadLocation("UTC")
	if err != nil {
		utcLoc = time.UTC
	}
}

// Reqn returns a unique transaction id (based on time) for WebMoney request
// Each subsequent request must have a Reqn greater than the previous one
func Reqn() string {
	locOnce.Do(initLocations)
	
	now := time.Now()
	nanoseconds := fmt.Sprintf("%03.f", float32(now.Nanosecond()/1000000))
	
	if os.Getenv("USE_OLD_REQN") != "" {
		return now.In(utcLoc).Format("060102150405") + nanoseconds[0:2]
	}
	return now.In(moscowLoc).Format("20060102150405") + nanoseconds
}

// Utf8ToWin encode string from utf8 to cp1251
func Utf8ToWin(s string) (string, error) {
	encoder := charmap.Windows1251.NewEncoder()
	return encoder.String(s)
}

func DebugLog(v ...interface{}) {
	if strings.ToUpper(os.Getenv("LOG_LEVEL")) == "TRACE" {
		var masked []interface{}
		for _, item := range v {
			if str, ok := item.(string); ok {
				masked = append(masked, MaskXMLSecrets(str))
			} else {
				masked = append(masked, item)
			}
		}
		fmt.Println(masked...)
	}
}

func DebugLogf(format string, v ...interface{}) {
	if strings.ToUpper(os.Getenv("LOG_LEVEL")) == "TRACE" {
		formatted := fmt.Sprintf(format, v...)
		fmt.Print(MaskXMLSecrets(formatted))
	}
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
