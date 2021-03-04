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
	return Request(context.Background(), "GET", url, options...)
}

func Post(url string, options ...Option) ([]byte, error) {
	return Request(context.Background(), "POST", url, options...)
}

func Request(ctx context.Context, method string, url string, options ...Option) ([]byte, error) {
	var opts = createOptions(options)

	request, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	// 重新配置request
	var values = opts.RequestBuilder(request)
	if values != "" {
		switch method {
		case "GET":
			request.URL.RawQuery = values
		case "POST":
			request.Body = ioutil.NopCloser(strings.NewReader(values))
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
