package cachex

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

/********************************************************************
created:    2021-08-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type CacheFuture struct {
	value      interface{}
	err        error
	updateTime unsafe.Pointer
	wg         sync.WaitGroup
}

func newCacheFuture() *CacheFuture {
	var item = &CacheFuture{}
	atomic.StorePointer(&item.updateTime, unsafe.Pointer(&time.Time{}))
	item.wg.Add(1)
	return item
}

func (my *CacheFuture) Get1() interface{} {
	my.wg.Wait()
	return my.value
}

func (my *CacheFuture) Get2() (interface{}, error) {
	my.wg.Wait()
	return my.value, my.err
}

// 这个方法只会被调用一次
func (my *CacheFuture) setValue(value interface{}, err error) {
	my.value = value
	my.err = err

	var now = time.Now()
	atomic.StorePointer(&my.updateTime, unsafe.Pointer(&now))
	my.wg.Done()
}

func (my *CacheFuture) getUpdateTime() time.Time {
	var p = (*time.Time)(atomic.LoadPointer(&my.updateTime))
	if p != nil {
		return *p
	}

	return time.Time{}
}
