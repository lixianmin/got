package loom

import (
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2020-09-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WheelTicker struct {
	wheel    *Wheel
	interval time.Duration
	data     *wheelData
}

func (my *WheelTicker) C() <-chan struct{} {
	if atomic.LoadInt32(&my.data.valid) == 0 {
		my.data = my.wheel.fetchWheelData(my.interval)
	}

	return my.data.c
}
