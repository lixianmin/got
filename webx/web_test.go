package webx

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

/********************************************************************
created:    2020-11-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestGet(t *testing.T) {
	var url = "http://www.baidu.com"
	var result, err = Get(url, WithTimeout(time.Second*2), WithRequestBuilder(func(request *http.Request) string {
		var query = request.URL.Query()
		query.Add("wd", "hello")
		var payload = query.Encode()

		return payload
	}))
	fmt.Printf("result=%s, err=%q", result, err)
}

func TestPostXWwwFormUrlencoded(t *testing.T) {
	var url = "http://www.baidu.com/s"
	var result, err = Post(url, WithRequestBuilder(func(request *http.Request) string {
		var header = request.Header
		header.Set("Content-Type", "application/x-www-form-urlencoded")

		var query = request.URL.Query()
		query.Add("wd", "hello")
		var payload = query.Encode()

		return payload
	}))

	fmt.Printf("result=%s, err=%q", result, err)
}

func TestPostJson(t *testing.T) {
	var url = "http://www.baidu.com"
	var result, err = Post(url, WithRequestBuilder(func(request *http.Request) string {
		var header = request.Header
		header.Set("Content-Type", "application/json")

		var payload = `{"name": "panda" }`
		return payload
	}))

	fmt.Printf("result=%s, err=%q", result, err)
}
