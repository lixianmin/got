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
	updateTime atomic.Value
	wg         sync.WaitGroup
	firstWait  sync.Once
}

func newCacheFuture() *CacheFuture {
	var item = &CacheFuture{}

	const year = 365 * 24 * time.Hour
	item.updateTime.Store(time.Now().Add(-year))
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
