package loom

import (
	"sync/atomic"
	"time"
	"unsafe"
)

/********************************************************************
created:    2020-09-11
author:     lixianmin

参考：https://blog.csdn.net/siddontang/article/details/18370541
https://github.com/siddontang/go/tree/master/timingwheel

Copyright (C) - All Rights Reserved
*********************************************************************/

type wheelData struct {
	c chan struct{}
}

type Wheel struct {
	wc          WaitClose
	step        time.Duration
	maxTimeout  time.Duration
	bucketsSize int
	position    int64
	channels    []unsafe.Pointer
}

func NewWheel(step time.Duration, bucketNum int) *Wheel {
	if step <= 0 {
		panic("step <= 0")
	}

	if bucketNum <= 0 {
		panic("bucketNum <= 0")
	}

	var wheel = &Wheel{
		step:        step,
		maxTimeout:  step * time.Duration(bucketNum),
		bucketsSize: bucketNum,
	}

	var channels = make([]unsafe.Pointer, bucketNum)
	for i := range channels {
		channels[i] = unsafe.Pointer(&wheelData{c: make(chan struct{})})
	}

	wheel.channels = channels
	Go(wheel.goLoop)
	return wheel
}

func (wheel *Wheel) Close() error {
	return wheel.wc.Close(nil)
}

func (wheel *Wheel) NewTimer(interval time.Duration) *WheelTimer {
	var timer = &WheelTimer{
		wheel:    wheel,
		interval: interval,
	}

	timer.reset(interval)
	return timer
}

func (wheel *Wheel) fetchWheelData(interval time.Duration) *wheelData {
	if interval < 0 || interval >= wheel.maxTimeout {
		panic("step should be in range [0, maxTimeout)")
	}

	var index = int(interval / wheel.step)
	if index > 0 {
		index--
	}

	var position = int(atomic.LoadInt64(&wheel.position))
	index = (position + index) % wheel.bucketsSize
	// 由于缺少lock控制，这里有可能取到已经被关闭的chan，但这没有关系，已经关闭的说明时刻已经过去了，立即返回就好
	var data = (*wheelData)(atomic.LoadPointer(&wheel.channels[index]))
	return data
}

func (wheel *Wheel) goLoop(later Later) {
	var ticker = later.NewTicker(wheel.step)
	var closeChan = wheel.wc.C()

	for {
		select {
		case <-ticker.C:
			wheel.onTicker()
		case <-closeChan:
			return
		}
	}
}

func (wheel *Wheel) onTicker() {
	var position = int(atomic.LoadInt64(&wheel.position))
	var lastItem = (*wheelData)(atomic.LoadPointer(&wheel.channels[position]))

	// 修改chan
	atomic.StorePointer(&wheel.channels[position], unsafe.Pointer(&wheelData{c: make(chan struct{})}))

	// 修改position
	position = (position + 1) % wheel.bucketsSize
	atomic.StoreInt64(&wheel.position, int64(position))

	// 关闭chan
	close(lastItem.c)
}
