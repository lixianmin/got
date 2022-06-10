package cachex

import (
	"runtime"
	"time"
)

/********************************************************************
created:    2022-04-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Loader = func(key interface{}) (interface{}, error)

type Cache interface {
	Load(key interface{}, loader Loader) *Future
	Set(key interface{}, value interface{}, err error)
	Get(key interface{}) interface{}
}

type wrapper struct {
	*cacheImpl
}

func NewCache(opts ...Option) Cache {
	var args = createArguments(opts)
	var my = &wrapper{&cacheImpl{
		args:      args,
		jobChan:   make(chan cacheJob, args.jobChanSize),
		gcTicker:  time.NewTicker(args.normalExpire * 4),
		closeChan: make(chan struct{}),
	}}

	// 初始化futures
	var shardingCount = cacheSharding.GetShardingCount()
	my.futures = make([]*cacheFuture, shardingCount)
	for i := 0; i < shardingCount; i++ {
		my.futures[i] = &cacheFuture{d: make(map[interface{}]*Future, 4)}
	}

	my.startJobGoroutines()
	// 参考: https://zhuanlan.zhihu.com/p/76504936
	runtime.SetFinalizer(my, func(w *wrapper) {
		w.gcTicker.Stop()
		close(w.closeChan) // 这里必须使用w.closeChan, 而不能使用my.closeChan, 否则runtime.GC()执行不到这里
		//fmt.Println("finalized")
	})
	return my
}
