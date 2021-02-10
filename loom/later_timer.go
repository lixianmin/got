package loom

import (
	"time"
)

/********************************************************************
created:    2021-02-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type LaterTimer struct {
	*time.Timer
	stoppedTime time.Time
}

func (my *LaterTimer) Stop() {
	my.Timer.Stop()
	my.stoppedTime = time.Time{}
}

func (my *LaterTimer) IsStopped() bool {
	return !time.Now().Before(my.stoppedTime)
}

// 这个方法用于支持随机时间区间的的Ticker()
// timer.Reset(timex.Duration(from, to))
func (my *LaterTimer) Reset(d time.Duration) bool {
	var b = my.Timer.Reset(d)
	my.stoppedTime = time.Now().Add(d)
	return b
}
