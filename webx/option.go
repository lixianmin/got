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
