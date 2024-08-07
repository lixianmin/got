package webx

import (
	"context"
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
	var url = "https://www.baidu.com"

	var ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var result, err = Get(ctx, url, WithRequestBuilder(func(request *http.Request) string {
		var query = request.URL.Query()
		query.Add("wd", "hello")
		var payload = query.Encode()

		return payload
	}))
	fmt.Printf("result=%s, err=%q", result, err)
}

func TestPostXWwwFormUrlencoded(t *testing.T) {
	var url = "http://www.baidu.com/s"
	var result, err = Post(context.Background(), url, WithRequestBuilder(func(request *http.Request) string {
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
	var result, err = Post(context.Background(), url, WithRequestBuilder(func(request *http.Request) string {
		var header = request.Header
		header.Set("Content-Type", "application/json")

		var payload = `{"name": "panda" }`
		return payload
	}))

	fmt.Printf("result=%s, err=%q", result, err)
}

//func TestGet2(t *testing.T) {
//	var url = "http://172.24.222.163:8888/hello"
//	var ctx, cancel = context.WithTimeout(context.Background(), 920*time.Millisecond)
//	defer cancel()
//
//	var result, err = Request(ctx, http.MethodGet, url)
//	fmt.Printf("result=%s, err=%q", result, err)
//}
