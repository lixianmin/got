package loom

import (
	"sync"
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2018-10-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WaitClose struct {
	C    chan struct{}
	m    sync.Mutex
	done int32
}

func NewWaitClose() *WaitClose {
	var wd = &WaitClose{
		C: make(chan struct{}),
	}

	return wd
}

func (wc *WaitClose) Close() {
	if 0 == atomic.LoadInt32(&wc.done) {
		wc.m.Lock()
		if 0 == wc.done {
			atomic.StoreInt32(&wc.done, 1)
			close(wc.C)
		}
		wc.m.Unlock()
	}
}

func (wc *WaitClose) WaitUtil(timeout time.Duration) bool {
	var timer = time.NewTimer(timeout)
	select {
	case <-wc.C:
		timer.Stop()
		return true
	case <-timer.C:
		return false
	}
}

func (wc *WaitClose) IsClosed() bool {
	return atomic.LoadInt32(&wc.done) == 1
}
