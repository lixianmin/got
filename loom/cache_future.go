package loom

import (
	"sync"
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2021-08-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type CacheFuture struct {
	value      interface{}
	err        error
	updateTime atomic.Value
	wg         sync.WaitGroup
}

func newCacheFuture() *CacheFuture {
	var item = &CacheFuture{}
	item.updateTime.Store(time.Time{})
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

	my.updateTime.Store(time.Now())
	my.wg.Done()
}

func (my *CacheFuture) getUpdateTime() time.Time {
	return my.updateTime.Load().(time.Time)
}