package loom

import (
	"time"
)

/********************************************************************
created:    2020-09-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WheelTimer struct {
	wheel    *Wheel
	interval time.Duration
	C        <-chan struct{}
}

// Reset()方法一定是在 <- timer.C 之后调用的，因此一定是在同一个 goroutine中，不存在竟态条件问题
func (my *WheelTimer) Reset() {
	var data = my.wheel.fetchWheelData(my.interval)
	my.C = data.c
}
