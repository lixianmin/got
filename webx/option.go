package webx

import (
	"net/http"
	"net/url"
	"time"
)

/********************************************************************
created:    2020-11-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type options struct {
	RequestBuilder func(request *http.Request) url.Values // 配置request
	Timeout        time.Duration                          // 控制从链接建立到返回的整个生命周期的时间
}

type Option func(*options)

func emptyRequestBuilder(request *http.Request) url.Values {
	return nil
}

func WithRequestBuilder(builder func(request *http.Request) url.Values) Option {
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
