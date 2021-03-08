package loom

import (
	"math/rand"
	"time"
)

/********************************************************************
created:    2018-10-12
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func Repeat(d time.Duration, handler func()) {
	go func() {
		defer DumpIfPanic()

		// 立即调用一次，保证及时初始化
		handler()

		// 在[d/2, d] 时间内调用一次，随机化ticker启动时间，防止缓存雪崩（同时出现大量瞬发性请求）
		var randomStart = time.Duration((int64(d) >> 1) + rand.Int63n(int64(int64(d)>>1)))
		time.Sleep(randomStart)

		for {
			handler()
			time.Sleep(d)
		}

		//// 立即调用一次，因为ticker需要过一段时间才触发
		//handler()
		//
		//// 在[0, d] 时间内调用一次，随机化ticker启动时间，防止大量的瞬发性请求
		//var randomStart = time.Duration(rand.Int63n(int64(d)))
		//time.Sleep(randomStart)
		//handler()
		//
		//// 使用ticker，每隔d的时间调用一次f()
		//var repeatTicker = time.NewTicker(d)
		//defer repeatTicker.Stop()
		//
		//for {
		//	select {
		//	case <-repeatTicker.C:
		//		handler()
		//	}
		//}
	}()
}