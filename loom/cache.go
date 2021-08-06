package loom

import (
	"sync"
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
	future *CacheFuture
}

type Cache struct {
	lockFutures sync.Mutex
	futures     map[interface{}]*CacheFuture
	jobChan     chan jobData
	wc          WaitClose
}

func NewCache(opts ...CacheOption) *Cache {
	var args = createCacheArguments(opts)
	var my = &Cache{
		futures: make(map[interface{}]*CacheFuture, 8),
		jobChan: make(chan jobData, 1),
	}

	my.startGoroutines(args.parallel)
	return my
}

func (my *Cache) startGoroutines(parallel int) {
	var jobChan = my.jobChan

	for i := 0; i < parallel; i++ {
		go func() {
			defer DumpIfPanic()
			for {
				select {
				case job := <-jobChan:
					var value = job.loader(job.key)
					job.future.setValue(value)
				case <-my.wc.C():
					break
				}
			}
		}()
	}
}

// Load 设计考量：
// 1. 如果缓存中有对应的Future对象，则直接返回
// 2. Load()方法自己不会阻塞，直接返回Future对象
// 3. 如果并发请求Load()方法，不会重复创建，会返回同一个Future对象
// 4. 超过2*expire的时间，则会移除超时的Future对象
// 5. 被移除的Future对象，如果已经被三方拿到了，可以正常调用Get()方法，如果内部正在加载，会正常加载完成
func (my *Cache) Load(key interface{}, expire time.Duration, loader func(key interface{}) interface{}) *CacheFuture {
	assert(key != nil, "key is nil")
	assert(loader != nil, "loader is nil")

	var future *CacheFuture
	var now = time.Now()
	my.lockFutures.Lock()
	{
		future = my.futures[key]
		// 如果future已经过期，则直接使用新的替换。否则如果返回旧future，用户可能Get()到过期的数据
		if future == nil || now.Sub(future.getUpdateTime()) >= 2*expire {
			future = newCacheFuture()
			my.futures[key] = future
		}
	}
	my.lockFutures.Unlock()

	var pastTime = now.Sub(future.getUpdateTime())
	if pastTime >= expire {
		my.jobChan <- jobData{
			loader: loader,
			key:    key,
			future: future,
		}
	}

	return future
}

func (my *Cache) Close() error {
	return my.wc.Close(nil)
}
