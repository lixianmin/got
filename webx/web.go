package webx

import (
	"context"
	"io"
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

	var request1, err1 = http.NewRequestWithContext(ctx, method, url, nil)
	if err1 != nil {
		return nil, err1
	}

	// 重新配置request
	var payload = opts.RequestBuilder(request1)
	if payload != "" {
		switch method {
		case http.MethodGet:
			request1.URL.RawQuery = payload
		case http.MethodPost:
			request1.Body = io.NopCloser(strings.NewReader(payload))
		}
	}

	var response2, err2 = opts.Client.Do(request1)
	if err2 != nil {
		return nil, err2
	}

	var responseBody = response2.Body
	defer responseBody.Close()

	var bts3, err3 = io.ReadAll(responseBody)
	return bts3, err3
}
