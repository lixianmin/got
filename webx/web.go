package webx

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/********************************************************************
created:    2018-09-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type RequestArgs struct {
	Builder func(request *http.Request) url.Values // 配置request
	Timeout time.Duration                          // 控制从链接建立到返回的整个生命周期的时间
}

func emptyRequestBuilder(request *http.Request) url.Values {
	return nil
}

func checkRequestArgs(args ...RequestArgs) RequestArgs {
	var arg = RequestArgs{}
	if len(args) > 0 {
		arg = args[0]
	}

	if arg.Timeout <= 0 {
		arg.Timeout = 10 * time.Second
	}

	if arg.Builder == nil {
		arg.Builder = emptyRequestBuilder
	}

	return arg
}

//func CopyHeader(dst, src http.Header) {
//	for k, vv := range src {
//		for _, v := range vv {
//			dst.Add(k, v)
//		}
//	}
//}

func Get(url string, args ...RequestArgs) ([]byte, error) {
	return Request(context.Background(), "GET", url, args...)
}

func Post(url string, args ...RequestArgs) ([]byte, error) {
	return Request(context.Background(), "POST", url, args...)
}

func Request(ctx context.Context, method string, url string, args ...RequestArgs) ([]byte, error) {
	var arg = checkRequestArgs(args...)

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// 重新配置request
	var values = arg.Builder(request)
	if values != nil {
		switch method {
		case "GET":
			request.URL.RawQuery = values.Encode()
		case "POST":
			request.Body = ioutil.NopCloser(strings.NewReader(values.Encode()))
		}
	}

	var client = http.Client{
		Timeout: arg.Timeout,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var responseBody = response.Body
	defer responseBody.Close()
	bodyBytes, err := ioutil.ReadAll(responseBody)
	return bodyBytes, err
}
