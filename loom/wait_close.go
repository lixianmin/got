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

// 2020-10-01
// 如果callback!=nil，则只要state != closed，就应该执行一下callback()
func (wc *WaitClose) Close(callback func() error) error {
	if wcClosed != atomic.LoadInt32(&wc.state) {
		wc.mutex.Lock()
		// 因为有外来的callback方法，所以有可能panic
		defer func() {
			wc.mutex.Unlock()
			if r := recover(); r != nil {
				fmt.Printf("%v\n", r)
			}
		}()

		// 双重检查
		if wcClosed != wc.state {
			if wcInitialized == wc.state {
				close(wc.closeChan)
			}

			// 即使未初始化的，也直接关闭掉
			atomic.StoreInt32(&wc.state, wcClosed)
			if callback != nil {
				return callback()
			}
		}
	}

	return nil
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
