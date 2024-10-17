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

func (wc *WaitClose) C() <-chan struct{} {
	if wcNew == atomic.LoadInt32(&wc.state) {
		wc.checkInitSlow()
	}

	return wc.closeChan
}

func (wc *WaitClose) WaitUtil(timeout time.Duration) bool {
	if wcNew == atomic.LoadInt32(&wc.state) {
		wc.checkInitSlow()
	}

	select {
	case <-wc.closeChan:
		return true
	case <-time.After(timeout):
		return false
	}
}

// Close 返回的时候，确保callback()方法已经执行完成
// 2020-10-01：如果callback!=nil，则只要state != closed，就应该执行一下callback()
// 2024-10-16: 如果callback中又调用了Close，则需要解决死锁问题
func (wc *WaitClose) Close(callback func() error) error {
	// 快速检查是否已关闭
	if wcClosed == atomic.LoadInt32(&wc.state) {
		return nil
	}

	var err error
	var needCallback bool

	// 第一阶段：尝试获取锁并更新状态
	wc.mutex.Lock()
	if wcClosed != wc.state {
		if wcInitialized == wc.state {
			close(wc.closeChan)
		} else {
			wc.closeChan = globalClosedChan
		}
		needCallback = true
		atomic.StoreInt32(&wc.state, wcClosed)
	}
	wc.mutex.Unlock()

	// 第二阶段：如果需要，在锁外执行回调
	if needCallback && callback != nil {
		err = safeCallback(callback)
	}

	return err
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

func safeCallback(callback func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in callback: %v\n", r)
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("panic in callback: %v", r)
			}
		}
	}()
	return callback()
}
