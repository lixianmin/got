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

var globalClosedChan = make(chan struct{})

func init() {
	// globalClosedChan用于保证wc.closeChan在任何情况下都不为nil，因为closed与nil这两种chan的表现很不一样
	close(globalClosedChan)
}

type WaitClose struct {
	closeChan chan struct{}
	mutex     sync.Mutex
	state     int32
}

func (wc *WaitClose) C() chan struct{} {
	if wcNew == atomic.LoadInt32(&wc.state) {
		wc.checkInitSlow()
	}
	//wc.assetCloseChanNotNil()

	return wc.closeChan
}

func (wc *WaitClose) WaitUtil(timeout time.Duration) bool {
	if wcNew == atomic.LoadInt32(&wc.state) {
		wc.checkInitSlow()
	}
	//wc.assetCloseChanNotNil()

	var timer = time.NewTimer(timeout)
	select {
	case <-wc.closeChan:
		timer.Stop()
		return true
	case <-timer.C:
		return false
	}
}

// Close 返回的时候，确保callback()方法已经执行完成
// 2020-10-01：如果callback!=nil，则只要state != closed，就应该执行一下callback()
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
			} else {
				wc.closeChan = globalClosedChan
			}

			// 即使未初始化的，也直接关闭掉
			// 2021-06-21 使用defer是为了确保Close()方法返回的时候，callback()方法必然已经执行完成
			defer atomic.StoreInt32(&wc.state, wcClosed)
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

// 1. 提取一个slow()方法，是为了让checkInit()方法可以inline
// 2. 进一步了，移除了checkInit()方法，手动inline了
func (wc *WaitClose) checkInitSlow() {
	wc.mutex.Lock()
	{
		if wcNew == wc.state {
			wc.closeChan = make(chan struct{})
			atomic.StoreInt32(&wc.state, wcInitialized)
		}
	}
	wc.mutex.Unlock()
}

func (wc *WaitClose) assetCloseChanNotNil() {
	// 下面这个断言有可能失败，好奇怪
	if wc.closeChan == nil {
		var message = fmt.Sprintf("closeChan is nil, state=%d, state=%d", wc.state, atomic.LoadInt32(&wc.state))
		panic(message)
	}
}
