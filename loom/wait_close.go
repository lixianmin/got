package loom

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2018-10-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const (
	wcNew = iota
	wcInitialized
	wcClosed
)

type WaitClose struct {
	closeChan chan struct{}
	mutex     sync.Mutex
	state     int32
}

func (wc *WaitClose) C() chan struct{} {
	wc.checkInit()
	return wc.closeChan
}

func (wc *WaitClose) WaitUtil(timeout time.Duration) bool {
	wc.checkInit()
	var timer = time.NewTimer(timeout)
	select {
	case <-wc.closeChan:
		timer.Stop()
		return true
	case <-timer.C:
		return false
	}
}

func (wc *WaitClose) Close(callback func()) {
	if wcInitialized == atomic.LoadInt32(&wc.state) {
		wc.mutex.Lock()
		if callback == nil {
			if wcInitialized == wc.state {
				atomic.StoreInt32(&wc.state, wcClosed)
				close(wc.closeChan)
			}
			wc.mutex.Unlock()
		} else { // 因为有外来的callback方法，所以有可能panic
			defer func() {
				wc.mutex.Unlock()
				if r := recover(); r != nil {
					fmt.Printf("%v\n", r)
				}
			}()

			if wcInitialized == wc.state {
				atomic.StoreInt32(&wc.state, wcClosed)
				close(wc.closeChan)
				callback()
			}
		}
	}
}

func (wc *WaitClose) IsClosed() bool {
	return atomic.LoadInt32(&wc.state) == wcClosed
}

func (wc *WaitClose) checkInit() {
	if wcNew == atomic.LoadInt32(&wc.state) {
		wc.mutex.Lock()
		if wcNew == wc.state {
			wc.closeChan = make(chan struct{})
			atomic.StoreInt32(&wc.state, wcInitialized)
		}
		wc.mutex.Unlock()
	}
}
