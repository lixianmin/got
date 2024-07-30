package webx

import (
	"net/http"
	"time"
)

/********************************************************************
created:    2020-11-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type RequestBuilderFunc func(request *http.Request) string
type options struct {
	RequestBuilder RequestBuilderFunc // 配置request
	Timeout        time.Duration      // 控制从链接建立到返回的整个生命周期的时间
}

type Option func(*options)

func emptyRequestBuilder(request *http.Request) string {
	return ""
}

// WithRequestBuilder
/*
 1. get方式编码:

	var result, err = webx.Get(context.Background(), url, webx.WithTimeout(time.Second*2), webx.WithRequestBuilder(func(request *http.Request) string {
		var query = request.URL.Query()
		query.Add("wd", "hello")
		var payload = query.Encode()

		return payload
	}))

 2. application/x-www-form-urlencoded方式编码:

	var result, err = webx.Post(context.Background(), url, webx.WithRequestBuilder(func(request *http.Request) string {
		var header = request.Header
		header.Set("Content-Type", "application/x-www-form-urlencoded")

		var query = request.URL.Query()
		query.Add("wd", "hello")
		var payload = query.Encode()

		return payload
	}))

 3.application/json编码方式:

	var result, err = webx.Post(context.Background(), url, webx.WithRequestBuilder(func(request *http.Request) string {
		var header = request.Header
		header.Set("Content-Type", "application/json")

		var payload = `{"name": "panda" }`
		return payload
	}))

*/
func WithRequestBuilder(builder RequestBuilderFunc) Option {
	return func(opts *options) {
		if builder != nil {
			opts.RequestBuilder = builder
		}
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(opts *options) {
		if timeout > 0 {
			opts.Timeout = timeout
		}
	}
}

func createOptions(optionList []Option) options {
	var opts = options{
		RequestBuilder: emptyRequestBuilder,
		Timeout:        10 * time.Second,
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	return opts
}
