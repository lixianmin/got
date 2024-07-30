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
	var opts = createOptions(options)

	// 讲道理, 无论是在http.Client{}中设置Timeout, 还是设置具备超时的ctx, 都是起作用的. 实测发现它们返回的err不一样
	// 现在遇到的问题是: 在设置http.Client{}的Timeout时, 请求仍然可能会超过很长时间, 因此现在尝试使用ctx的方案, 看看是否可以避免这个问题
	// 实测不解决问题, 还是会超时
	//var ctx1, cancel = context.WithTimeout(ctx, opts.Timeout)
	//defer cancel()

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
