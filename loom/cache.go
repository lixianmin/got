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

type loadJob struct {
	loader func(key interface{}) interface{}
	key    interface{}
	future *CacheFuture
}

type Cache struct {
	expire      time.Duration
	lockFutures sync.Mutex
	futures     map[interface{}]*CacheFuture
	lockJob     sync.Mutex
	jobChan     chan loadJob
	wc          WaitClose
}

func NewCache(opts ...CacheOption) *Cache {
	var args = createCacheArguments(opts)
	var my = &Cache{
		expire:  args.expire,
		futures: make(map[interface{}]*CacheFuture, 8),
		jobChan: make(chan loadJob, 1),
	}

	my.startGoroutines(args.parallel)
	Repeat(time.Minute, my.removeExpired)
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
					job.future.setLoading(false)
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
func (my *Cache) Load(key interface{}, loader func(key interface{}) interface{}) *CacheFuture {
	assert(key != nil, "key is nil")
	assert(loader != nil, "loader is nil")

	var future = my.fetchFuture(key)
	var needUpdate = time.Now().Sub(future.getUpdateTime()) > my.expire
	if needUpdate {
		my.checkStartLoad(future, key, loader)
	}

	return future
}

func (my *Cache) fetchFuture(key interface{}) *CacheFuture {
	var future *CacheFuture
	my.lockFutures.Lock()
	{
		future = my.futures[key]
		// 如果future已经过期，则直接使用新的替换。否则如果返回旧future，用户可能Get()到过期的数据
		if future == nil || my.shouldRemove(future) {
			future = newCacheFuture()
			my.futures[key] = future
		}
	}
	my.lockFutures.Unlock()

	return future
}

func (my *Cache) checkStartLoad(future *CacheFuture, key interface{}, loader func(key interface{}) interface{}) {
	my.lockJob.Lock()
	{
		if !future.isLoading() {
			future.setLoading(true)
			my.jobChan <- loadJob{
				loader: loader,
				key:    key,
				future: future,
			}
		}
	}
	my.lockJob.Unlock()
}

func (my *Cache) Close() error {
	return my.wc.Close(nil)
}

func (my *Cache) removeExpired() {
	my.lockFutures.Lock()
	{
		for key, future := range my.futures {
			if my.shouldRemove(future) {
				delete(my.futures, key)
			}
		}
	}
	my.lockFutures.Unlock()
}

func (my *Cache) shouldRemove(future *CacheFuture) bool {
	var updateTime = future.getUpdateTime()
	return updateTime != time.Time{} && time.Now().Sub(updateTime) > 2*my.expire
}
