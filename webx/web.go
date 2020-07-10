package webx

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

/********************************************************************
created:    2018-09-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type RequestArgs struct {
	RequestBuilder func(request *http.Request) // 初始化request
	Timeout        time.Duration               // 控制从链接建立到返回的整个生命周期的时间
}

func emptyRequestBuilder(request *http.Request) {

}

func checkRequestArgs(args *RequestArgs) *RequestArgs {
	if nil == args {
		args = &RequestArgs{}
	}

	if args.Timeout <= 0 {
		args.Timeout = 10 * time.Second
	}

	if args.RequestBuilder == nil {
		args.RequestBuilder = emptyRequestBuilder
	}

	return args
}

//func CopyHeader(dst, src http.Header) {
//	for k, vv := range src {
//		for _, v := range vv {
//			dst.Add(k, v)
//		}
//	}
//}

func Get(url string, args *RequestArgs) ([]byte, error) {
	return Request(context.Background(), "GET", url, args)
}

func Post(url string, args *RequestArgs) ([]byte, error) {
	return Request(context.Background(), "POST", url, args)
}

func Request(ctx context.Context, method string, url string, args *RequestArgs) ([]byte, error) {
	args = checkRequestArgs(args)

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// 重新配置request
	args.RequestBuilder(request)

	var client = http.Client{
		Timeout: args.Timeout,
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
