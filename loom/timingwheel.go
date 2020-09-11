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

type wheelChan struct {
	c chan struct{}
}

type TimingWheel struct {
	wc          WaitClose
	interval    time.Duration
	maxTimeout  time.Duration
	bucketsSize int
	position    int64
	channels    []unsafe.Pointer
}

func NewTimingWheel(interval time.Duration, bucketsSize int) *TimingWheel {
	if interval <= 0 {
		panic("interval <= 0")
	}

	if bucketsSize <= 0 {
		panic("bucketsSize <= 0")
	}

	var wheel = &TimingWheel{
		interval:    interval,
		maxTimeout:  interval * time.Duration(bucketsSize),
		bucketsSize: bucketsSize,
	}

	var channels = make([]unsafe.Pointer, bucketsSize)
	for i := range channels {
		channels[i] = unsafe.Pointer(&wheelChan{c: make(chan struct{})})
	}

	wheel.channels = channels
	Go(wheel.goLoop)
	return wheel
}

func (wheel *TimingWheel) Close() error {
	return wheel.wc.Close(nil)
}

func (wheel *TimingWheel) After(timeout time.Duration) <-chan struct{} {
	if timeout >= wheel.maxTimeout {
		panic("timeout too much, over max timeout")
	}

	var index = int(timeout / wheel.interval)
	if index > 0 {
		index--
	}

	var position = int(atomic.LoadInt64(&wheel.position))
	index = (position + index) % wheel.bucketsSize
	var waitChan = (*wheelChan)(atomic.LoadPointer(&wheel.channels[index])).c

	return waitChan
}

func (wheel *TimingWheel) goLoop(later Later) {
	var ticker = later.NewTicker(wheel.interval)
	var closeChan = wheel.wc.C()

	for {
		select {
		case <-ticker.C:
			wheel.onTicker()
		case <-closeChan:
			ticker.Stop()
			return
		}
	}
}

func (wheel *TimingWheel) onTicker() {
	var position = int(atomic.LoadInt64(&wheel.position))
	var lastChan = (*wheelChan)(atomic.LoadPointer(&wheel.channels[position])).c
	atomic.StorePointer(&wheel.channels[position], unsafe.Pointer(&wheelChan{c: make(chan struct{})}))

	position = (position + 1) % wheel.bucketsSize
	atomic.StoreInt64(&wheel.position, int64(position))
	close(lastChan)
}
