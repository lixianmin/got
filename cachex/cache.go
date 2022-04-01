package cachex

import (
	"fmt"
	"runtime"
	"time"
)

/********************************************************************
created:    2022-04-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type CacheLoader = func(key interface{}) (interface{}, error)

//type Cache interface {
//	Load(key interface{}, loader CacheLoader) *CacheFuture
//}

type Cache struct {
	*CacheImpl
}

func NewCache(opts ...CacheOption) *Cache {
	var args = createArguments(opts)
	var my = &Cache{&CacheImpl{
		args:      args,
		jobChan:   make(chan cacheJob, args.jobChanSize),
		gcTicker:  time.NewTicker(args.normalExpire * 4),
		closeChan: make(chan struct{}),
	}}

	// 初始化futures
	var shardingCount = cacheSharding.GetShardingCount()
	my.futures = make([]*cacheFuture, shardingCount)
	for i := 0; i < shardingCount; i++ {
		my.futures[i] = &cacheFuture{d: make(map[interface{}]*CacheFuture, 4)}
	}

	my.startJobGoroutines()
	runtime.SetFinalizer(my, func(my *Cache) {
		my.gcTicker.Stop()
		close(my.closeChan)
		fmt.Println("finalized")
	})
	return my
}
