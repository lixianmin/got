package webx

import (
	"io/ioutil"
	"net/http"
	"time"
)

/********************************************************************
created:    2018-09-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func CopyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func Get(url string, initHeader func(header http.Header)) ([]byte, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if initHeader != nil {
		initHeader(request.Header)
	}

	var client = http.Client{
		Timeout: time.Second * 20, // 控制从链接建立到返回的整个生命周期的时间
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	var body = response.Body
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	return bodyBytes, err
}
