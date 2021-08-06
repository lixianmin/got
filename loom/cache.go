package loom

import (
	"time"
)

/********************************************************************
created:    2021-08-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type jobData struct {
	loader func(key interface{}) interface{}
	key    interface{}
	item   *CacheItem
}

type Cache struct {
	items   Map
	jobChan chan jobData
	wc      WaitClose
}

func NewCache(parallel int) *Cache {
	if parallel < 1 {
		panic("parallel is to small")
	}

	var my = &Cache{
		jobChan: make(chan jobData, 1),
	}

	my.startGoroutines(parallel)
	return my
}

func (my *Cache) startGoroutines(parallel int) {
	for i := 0; i < parallel; i++ {
		go func() {
			defer DumpIfPanic()
			for {
				select {
				case job := <-my.jobChan:
					var value = job.loader(job.key)
					job.item.setValue(value)
				case <-my.wc.C():
					break
				}
			}
		}()
	}
}

func (my *Cache) Load(key interface{}, expire time.Duration, loader func(key interface{}) interface{}) *CacheItem {
	var item = my.items.ComputeIfAbsent(key, func(key interface{}) interface{} {
		return newCacheItem()
	}).(*CacheItem)

	var updateTime = item.getUpdateTime()
	var pastTime = time.Now().Sub(updateTime)
	if pastTime >= expire {
		my.jobChan <- jobData{
			loader: loader,
			key:    key,
			item:   item,
		}
	}

	return item
}

func (my *Cache) Close() error {
	return my.wc.Close(nil)
}
