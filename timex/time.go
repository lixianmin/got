package timex

import (
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2020-07-23
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const TimeLayout = "2006-01-02 15:04:05"

var currentTime = time.Now().Unix()

func init() {
	go func() {
		for {
			atomic.StoreInt64(&currentTime, time.Now().Unix())
			time.Sleep(time.Second)
		}
	}()
}

// 按照国人习惯的方式格式化了一下时间
func FormatTime(t time.Time) string {
	return t.Format(TimeLayout)
}

// 返回一个秒级的低精度的时间戳
func NowUnix() int64 {
	return atomic.LoadInt64(&currentTime)
}
