// Copyright 2015 Constantin Karataev. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// This file contains some helper function to form data for requests

package webmoney

import (
	"log"
	//"strconv"
	"time"
	"fmt"

	"github.com/qiniu/iconv"
)

// function return reqn params for webmoney request
// each subsequent more then the previews
func Reqn() string {
	nanoseconds := fmt.Sprintf("%03.f",float32(time.Now().Nanosecond()/1000000))
	return time.Now().Local().Format("20060102150405") + nanoseconds
}

// encode string from utf8 to cp1251
func Utf8ToWin(s string) (string, error) {
	d, err := iconv.Open("cp1251//IGNORE", "utf8")
	if err != nil {
		return "", err
	}
	defer d.Close()
	return d.ConvString(s), nil
}

func Str4Sign(s string) string {
	res, err := Utf8ToWin(s)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return res
}

// entrode string from cp1251 to utf8
func WinToUtf8(s string) (string, error) {
	d, err := iconv.Open("utf8", "cp1251")
	if err != nil {
		return "", err
	}
	defer d.Close()
	return d.ConvString(s), nil
}
