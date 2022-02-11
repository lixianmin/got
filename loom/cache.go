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

type CacheLoader = func(key interface{}) (interface{}, error)

type cacheJob struct {
	loader CacheLoader
	key    interface{}
	future *CacheFuture
}

type cacheFuture struct {
	sync.RWMutex
	d map[interface{}]*CacheFuture
}

type Cache struct {
	expire   time.Duration
	futures  []*cacheFuture
	lockJob  sync.Mutex
	jobChan  chan cacheJob
	gcTicker *time.Ticker
	wc       WaitClose
}

func NewCache(opts ...CacheOption) *Cache {
	var args = createCacheArguments(opts)
	var my = &Cache{
		expire:   args.expire,
		jobChan:  make(chan cacheJob, 128), // 加大这个chan的长度, 有助于减小第一次checkLoad()时的执行时间
		gcTicker: time.NewTicker(args.gcInterval),
	}

	// 初始化futures
	var shardingCount = mapSharding.GetShardingCount()
	my.futures = make([]*cacheFuture, shardingCount)
	for i := 0; i < shardingCount; i++ {
		my.futures[i] = &cacheFuture{d: make(map[interface{}]*CacheFuture, 4)}
	}

	my.startGoroutines(args)
	return my
}

func (my *Cache) startGoroutines(args cacheArguments) {
	var jobChan = my.jobChan

	for i := 0; i < args.parallel; i++ {
		go func() {
			defer DumpIfPanic()
			for {
				select {
				case job := <-jobChan:
					var value, err = job.loader(job.key)
					job.future.setValue(value, err)
				case <-my.gcTicker.C:
					my.removeRotted()
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
// 4. 超过2*expire的时间，称之为rotted，会直接移除Future对象
// 5. 被移除的Future对象，如果已经被三方拿到了，可以正常调用Get()方法，如果内部正在加载，会正常加载完成
func (my *Cache) Load(key interface{}, loader CacheLoader) *CacheFuture {
	assert(key != nil, "key is nil")
	assert(loader != nil, "loader is nil")

	var future = my.fetchFuture(key)
	var mayNeedLoad = time.Now().Sub(future.getUpdateTime()) > my.expire
	if mayNeedLoad {
		my.checkLoad(future, key, loader)
	}

	return future
}

//func (my *Cache) fetchFuture1(key interface{}) *CacheFuture {
//	var future *CacheFuture
//	var futures = my.futures
//	futures.Lock()
//	{
//		future = futures.d[key]
//		// 如果future已经rotted，则直接使用新的替换。否则如果返回旧future，用户可能Get()到过期的数据
//		if future == nil || my.isRotted(future) {
//			future = newCacheFuture()
//			futures.d[key] = future
//		}
//	}
//	futures.Unlock()
//
//	return future
//}

// 使用RWMutex比Mutex速度提高1/3左右
func (my *Cache) fetchFuture(key interface{}) *CacheFuture {
	var index, _ = mapSharding.GetShardingIndex(key)
	var futures = my.futures[index]
	var future *CacheFuture

	// 尝试获取缓存中的future, 如果已经rotted, 则不返回它
	futures.RLock()
	{
		future = futures.d[key]
		if future != nil && my.isRotted(future) {
			future = nil
		}
	}
	futures.RUnlock()

	// 如果future为nil, 则代表缓存中不存在或者已经rotted
	if future == nil {
		futures.Lock()
		{
			future = newCacheFuture()
			futures.d[key] = future
		}
		futures.Unlock()
	}

	return future
}

func (my *Cache) checkLoad(future *CacheFuture, key interface{}, loader CacheLoader) {
	// fast path
	if future.getStatus() == kFutureInit {
		my.checkLoadSlowPath(future, key, loader)
	}
}

func (my *Cache) checkLoadSlowPath(future *CacheFuture, key interface{}, loader CacheLoader) {
	// lockJob这把锁, 放到Cache而不是CacheFuture中的原因是:
	//  1. 节约内存
	//  2. benchmark测试性能区别不大, 估计是因为fast path的原因被均摊了
	my.lockJob.Lock()
	{
		if future.getStatus() == kFutureInit {
			future.setStatus(kFutureLoading)
			my.jobChan <- cacheJob{
				loader: loader,
				key:    key,
				future: future,
			}
		}
	}
	my.lockJob.Unlock()
}

func (my *Cache) Close() error {
	return my.wc.Close(func() error {
		my.gcTicker.Stop()
		return nil
	})
}

func (my *Cache) removeRotted() {
	for _, futures := range my.futures {
		futures.Lock()
		{
			for key, future := range futures.d {
				if my.isRotted(future) {
					delete(futures.d, key)
				}
			}
		}
		futures.Unlock()
	}
}

//func (my *Cache) isExpired(future *CacheFuture) bool {
//	var updateTime = future.getUpdateTime()
//	return updateTime != time.Time{} && time.Now().Sub(updateTime) > my.expire
//}

// 超过1倍的expire，称为『过期』
// 超过2倍的expire，称为『腐烂』
func (my *Cache) isRotted(future *CacheFuture) bool {
	var updateTime = future.getUpdateTime()
	return updateTime != time.Time{} && time.Now().Sub(updateTime) > 2*my.expire
}
