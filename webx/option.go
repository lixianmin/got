package webx

import (
	"net/http"
)

/********************************************************************
created:    2020-11-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type RequestBuilderFunc func(request *http.Request) string
type options struct {
	RequestBuilder RequestBuilderFunc // 配置request
	Client         *http.Client       // http.Client中有链接池 (可通过http.Transport配置), 因此不宜于每次请求生成一个http.Client对象
}

type Option func(*options)

func emptyRequestBuilder(request *http.Request) string {
	return ""
}

// WithRequestBuilder
/*
 1. get方式编码:

	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var result, err = webx.Get(ctx, url, webx.WithRequestBuilder(func(request *http.Request) string {
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

func WithClient(client *http.Client) Option {
	return func(opts *options) {
		if client != nil {
			opts.Client = client
		}
	}
}

func createOptions(optionList []Option) options {
	var opts = options{
		RequestBuilder: emptyRequestBuilder,
		Client:         http.DefaultClient, // http.DefaultClient中没有设置Timeout
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	// // 讲道理, 无论是在http.Client{}中设置Timeout, 还是设置具备超时的ctx, 都是起作用的. 实测发现它们返回的err不一样
	// // 如果ctx中设置了timeout, 则也使用这个timeout设置到http身上
	// var deadline, ok = ctx.Deadline()
	// if ok {
	// 	opts.Timeout = time.Until(deadline)
	// }

	return opts
}
