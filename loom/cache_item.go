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

type CacheItem struct {
	value      interface{}
	updateTime atomic.Value
	wg         sync.WaitGroup
	firstWait  sync.Once
}

func newCacheItem() *CacheItem {
	var item = &CacheItem{}

	const year = 365 * 24 * time.Hour
	item.updateTime.Store(time.Now().Add(-year))
	item.wg.Add(1)
	return item
}

func (my *CacheItem) Get() interface{} {
	my.wg.Wait()
	return my.value
}

func (my *CacheItem) setValue(value interface{}) {
	my.value = value
	my.updateTime.Store(time.Now())
	my.firstWait.Do(func() {
		my.wg.Done()
	})
}

func (my *CacheItem) getUpdateTime() time.Time {
	return my.updateTime.Load().(time.Time)
}
