package loom

import (
	"sync"
	"time"
)

/********************************************************************
created:    2020-09-11
author:     lixianmin

参考：https://blog.csdn.net/siddontang/article/details/18370541
https://github.com/siddontang/go/tree/master/timingwheel

Copyright (C) - All Rights Reserved
*********************************************************************/

type TimingWheel struct {
	wc WaitClose

	interval    time.Duration
	maxTimeout  time.Duration
	bucketsSize int

	buckets struct {
		sync.RWMutex
		channels []chan struct{}
		position int
	}
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

	var channels = make([]chan struct{}, bucketsSize)
	for i := range channels {
		channels[i] = make(chan struct{})
	}

	wheel.buckets.channels = channels
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

	var buckets = &wheel.buckets
	buckets.RLock()
	index = (buckets.position + index) % wheel.bucketsSize
	var waitChan = buckets.channels[index]
	buckets.RUnlock()

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
	var buckets = &wheel.buckets
	var nextChan = make(chan struct{})

	buckets.Lock()
	var position = buckets.position
	var lastChan = buckets.channels[position]
	buckets.channels[position] = nextChan
	buckets.position = (position + 1) % wheel.bucketsSize
	buckets.Unlock()

	close(lastChan)
}
