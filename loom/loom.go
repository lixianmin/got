package loom

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime/debug"
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

func DumpIfPanic() {
	var panicData = recover()
	if panicData == nil {
		return
	}

	var exeName = filepath.Base(os.Args[0]) // 获取程序名称

	// 设定时间格式
	const format = "2006-01-02T15:04:05"
	var timestamp = time.Now().Format(format)
	// 保存错误信息文件名:dump.程序名.当前时间（年月日时分秒）
	var logDir = "logs"
	var logFilePath = fmt.Sprintf("%s/dump.%s.%s.log", logDir, exeName, timestamp)
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		_ = os.MkdirAll(logDir, os.ModePerm)
	}

	fmt.Println("dump to file ", logFilePath)

	f, err := os.Create(logFilePath)
	if err != nil {
		return
	}
	defer f.Close()

	// 输出panic信息
	//writeOneMessage(f, "------------------------------------\r\n")
	//writeOneMessage(f, message)
	writeOneMessage(f, "------------------------------------\r\n")
	writeOneMessage(f, fmt.Sprintf("%v\r\n", panicData))
	writeOneMessage(f, "------------------------------------\r\n")

	// 输出堆栈信息
	writeOneMessage(f, string(debug.Stack()))

	// 直接退出？
	os.Exit(1)
}

func writeOneMessage(fout *os.File, message string) {
	_, _ = fout.WriteString(message)
	_, _ = os.Stderr.WriteString(message)
}
