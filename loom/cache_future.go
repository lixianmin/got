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
	firstWait  sync.Once
	loading    int32
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

func (my *CacheFuture) setValue(value interface{}, err error) {
	my.value = value
	my.err = err

	my.updateTime.Store(time.Now())
	my.firstWait.Do(func() {
		my.wg.Done()
	})
}

func (my *CacheFuture) getUpdateTime() time.Time {
	return my.updateTime.Load().(time.Time)
}

func (my *CacheFuture) setLoading(b bool) {
	var loading int32 = 0
	if b {
		loading = 1
	}
	atomic.StoreInt32(&my.loading, loading)
}

func (my *CacheFuture) isLoading() bool {
	var d = atomic.LoadInt32(&my.loading)
	return d == 1
}
