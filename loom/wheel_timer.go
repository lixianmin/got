package loom

import (
	"sync/atomic"
	"time"
	"unsafe"
)

/********************************************************************
created:    2020-09-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WheelTimer struct {
	wheel    *Wheel
	interval time.Duration
	data     unsafe.Pointer
}

func (my *WheelTimer) Restart() {
	my.Reset(my.interval)
}

func (my *WheelTimer) Reset(interval time.Duration) {
	var data = my.wheel.fetchWheelData(interval)
	atomic.StorePointer(&my.data, unsafe.Pointer(data))
}

func (my *WheelTimer) C() <-chan struct{} {
	var data = (*wheelData)(atomic.LoadPointer(&my.data))
	return data.c
}
