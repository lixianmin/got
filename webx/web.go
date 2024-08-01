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

func Get(ctx context.Context, url string, options ...Option) ([]byte, error) {
	return Request(ctx, http.MethodGet, url, options...)
}

func Post(ctx context.Context, url string, options ...Option) ([]byte, error) {
	return Request(ctx, http.MethodPost, url, options...)
}

func Request(ctx context.Context, method string, url string, options ...Option) ([]byte, error) {
	var opts = createOptions(ctx, options)

	request, err1 := http.NewRequestWithContext(ctx, method, url, nil)
	if err1 != nil {
		return nil, err1
	}

	// 重新配置request
	var payload = opts.RequestBuilder(request)
	if payload != "" {
		switch method {
		case http.MethodGet:
			request.URL.RawQuery = payload
		case http.MethodPost:
			request.Body = ioutil.NopCloser(strings.NewReader(payload))
		}
	}

	var client = http.Client{
		Timeout: opts.Timeout,
	}

	response, err2 := client.Do(request)
	if err2 != nil {
		return nil, err2
	}

	var responseBody = response.Body
	defer responseBody.Close()

	bodyBytes, err3 := ioutil.ReadAll(responseBody)
	return bodyBytes, err3
}
