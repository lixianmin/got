package loom

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"
)

/********************************************************************
created:    2020-03-08
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var dumpHandler func(data []byte)

func Initialize(dumpCallback func(data []byte)) {
	dumpHandler = dumpCallback
}

func DumpIfPanic() {
	var panicData = recover()
	if panicData == nil {
		return
	}

	var serviceName = filepath.Base(os.Args[0]) // 获取程序名称

	// 设定时间格式
	const format = "2006-01-02T15:04:05"
	var timestamp = time.Now().Format(format)
	// 保存错误信息文件名:dump.程序名.当前时间（年月日时分秒）
	var logDir = "logs"
	var logFilePath = fmt.Sprintf("%s/dump.%s.%s.log", logDir, serviceName, timestamp)
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
	const lineSeparator = "------------------------------------\r\n"
	var data = make([]byte, 0, 1024)
	data = append(data, lineSeparator...)
	data = append(data, fmt.Sprintf("%v\r\n", panicData)...)
	data = append(data, lineSeparator...)
	data = append(data, debug.Stack()...) // 调用栈信息

	// 输出
	_, _ = f.Write(data)
	if nil != dumpHandler {
		dumpHandler(data)
	}

	// 直接退出？
	os.Exit(1)
}
