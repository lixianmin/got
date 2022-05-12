package webx

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
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

func Get(url string, options ...Option) ([]byte, error) {
	return Request(context.Background(), http.MethodGet, url, options...)
}

func Post(url string, options ...Option) ([]byte, error) {
	return Request(context.Background(), http.MethodPost, url, options...)
}

func Request(ctx context.Context, method string, url string, options ...Option) ([]byte, error) {
	var opts = createOptions(options)

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// 重新配置request
	var payload = opts.RequestBuilder(request)
	if payload != "" {
		switch method {
		case "GET":
			request.URL.RawQuery = payload
		case "POST":
			request.Body = ioutil.NopCloser(strings.NewReader(payload))
		}
	}

	var client = http.Client{
		Timeout: opts.Timeout,
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
