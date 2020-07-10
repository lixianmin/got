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

//func CopyHeader(dst, src http.Header) {
//	for k, vv := range src {
//		for _, v := range vv {
//			dst.Add(k, v)
//		}
//	}
//}

func Get(url string, initRequest func(request *http.Request)) ([]byte, error) {
	return Request(context.Background(), "GET", url, initRequest)
}

func Post(url string, initRequest func(request *http.Request)) ([]byte, error) {
	return Request(context.Background(), "POST", url, initRequest)
}

func Request(ctx context.Context, method string, url string, initRequest func(request *http.Request)) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	if initRequest != nil {
		initRequest(request)
	}

	var client = http.Client{
		Timeout: time.Second * 10, // 控制从链接建立到返回的整个生命周期的时间
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
