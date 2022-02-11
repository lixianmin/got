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

const (
	kFutureEmpty   = iota // 无, 则new&load, 返回new
	kFutureGood           // 有 & 可用, 返回old
	kFutureExpired        // 有 & 过期, new&load, 返回old
	kFutureRotted         // 有 & rotted, new&load, 返回new
)

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
	expire  time.Duration
	futures []*cacheFuture
	jobChan chan cacheJob
	wc      WaitClose
}

func NewCache(opts ...CacheOption) *Cache {
	var args = createCacheArguments(opts)
	var my = &Cache{
		expire:  args.expire,
		jobChan: make(chan cacheJob, 128), // 加大这个chan的长度, 有助于减小第一次checkLoad()时的执行时间
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
	var closeChan = my.wc.C()

	var ticker = time.NewTicker(args.gcInterval)
	var stopOnce sync.Once

	for i := 0; i < args.parallel; i++ {
		go func() {
			defer DumpIfPanic()
			defer stopOnce.Do(func() {
				ticker.Stop()
			})

			for {
				select {
				case job := <-jobChan:
					var value, err = job.loader(job.key)
					job.future.setValue(value, err)
				case <-ticker.C:
					my.removeRotted()
				case <-closeChan:
					return
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

	var index, _ = mapSharding.GetShardingIndex(key)
	var futures = my.futures[index]

	// 尝试获取缓存中的future
	futures.RLock()
	var last = futures.d[key]
	futures.RUnlock()

	var status = my.getFutureStatus(last)
	if status == kFutureGood {
		return last
	}

	var next = newCacheFuture()
	futures.Lock()
	futures.d[key] = next
	futures.Unlock()
	my.jobChan <- cacheJob{loader: loader, key: key, future: next}

	if status == kFutureExpired {
		return last
	}

	return next
}

func (my *Cache) Close() error {
	return my.wc.Close(nil)
}

func (my *Cache) removeRotted() {
	for _, futures := range my.futures {
		futures.Lock()
		{
			for key, future := range futures.d {
				if my.getFutureStatus(future) == kFutureRotted {
					delete(futures.d, key)
				}
			}
		}
		futures.Unlock()
	}
}

func (my *Cache) getFutureStatus(last *CacheFuture) int {
	if last != nil {
		var updateTime = last.getUpdateTime()
		var past = time.Now().Sub(updateTime)
		var expire = my.expire

		if updateTime.IsZero() || past < expire {
			return kFutureGood
		} else if past < 2*expire {
			return kFutureExpired
		} else {
			return kFutureRotted
		}
	}

	return kFutureEmpty
}
