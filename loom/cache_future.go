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
	value          interface{}
	expireDuration int64
	updateTime     atomic.Value
	wg             sync.WaitGroup
	firstWait      sync.Once
	loading        int32
}

func newCacheFuture() *CacheFuture {
	var item = &CacheFuture{}
	item.updateTime.Store(time.Time{})
	item.wg.Add(1)
	return item
}

func (my *CacheFuture) Get() interface{} {
	my.wg.Wait()
	return my.value
}

func (my *CacheFuture) setValue(value interface{}) {
	my.value = value
	my.updateTime.Store(time.Now())
	my.firstWait.Do(func() {
		my.wg.Done()
	})
}

func (my *CacheFuture) getUpdateTime() time.Time {
	return my.updateTime.Load().(time.Time)
}

func (my *CacheFuture) setExpireDuration(expire time.Duration) {
	atomic.StoreInt64(&my.expireDuration, int64(expire))
}

func (my *CacheFuture) getExpireDuration() time.Duration {
	return time.Duration(atomic.LoadInt64(&my.expireDuration))
}

func (my *CacheFuture) isExpired(expire time.Duration) bool {
	var updateTime = my.getUpdateTime()
	return updateTime != time.Time{} && time.Now().Sub(updateTime) > expire
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
