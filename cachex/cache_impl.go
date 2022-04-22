package cachex

import (
	"github.com/lixianmin/got/loom"
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
	kFutureGood           // 有 & 可用, 返回last
	kFutureExpired        // 有 & 过期, new & load, 返回 last
	kFutureRotted         // 有 & rotted, new & load, 返回 new
)

var cacheSharding = loom.NewSharding()

type cacheJob struct {
	loader Loader
	key    interface{}
	future *Future
}

type cacheFuture struct {
	sync.Mutex
	d map[interface{}]*Future
}

type cacheImpl struct {
	args      arguments
	futures   []*cacheFuture
	jobChan   chan cacheJob
	gcTicker  *time.Ticker
	closeChan chan struct{}
}

func (my *cacheImpl) startJobGoroutines() {
	var jobChan = my.jobChan
	var parallel = my.args.parallel
	var gcTicker = my.gcTicker
	var closeChan = my.closeChan

	for i := 0; i < parallel; i++ {
		go func() {
			defer loom.DumpIfPanic()

			for {
				select {
				case job := <-jobChan:
					var value, err = job.loader(job.key)
					job.future.setValue(value, err)
				case <-gcTicker.C:
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
// 4. 被移除的Future对象，如果已经被三方拿到了，可以正常调用Get()方法，如果内部正在加载，会正常加载完成
func (my *cacheImpl) Load(key interface{}, loader Loader) *Future {
	assert(key != nil, "key is nil")
	assert(loader != nil, "loader is nil")

	var index, _ = cacheSharding.GetShardingIndex(key)
	var futures = my.futures[index]
	var next *Future = nil

	// 以下代码需要考虑并发, 需要阻止重复加载
	futures.Lock()
	var last = futures.d[key]
	var lastStatus = my.getFutureStatus(last)
	// todo 如果是expired的future, 本次会返回last, 但立马会把next写入, 下次就会是good状态但未加载完的对象了

	// 可能是empty, expired, rotted
	if lastStatus != kFutureGood {
		var predecessor *Future = nil
		if lastStatus == kFutureExpired { // 只有处于expired状态的last才有资格当作predecessor
			predecessor = last
		}

		next = newFuture(predecessor)
		futures.d[key] = next
		my.sendJob(cacheJob{loader: loader, key: key, future: next})
	}
	futures.Unlock()

	//fmt.Printf("lastStatus=%v \n", lastStatus)
	switch lastStatus {
	case kFutureGood: // lastStatus == good: 意味着last本身还没有加载完呢, 所以不会创建next, 因此不可能返回next. 但是, 这个last有可能有可勉强使用的predecessor
		var predecessor = last.getPredecessor()
		var status = my.getFutureStatus(predecessor)
		//fmt.Printf("status=%d \n", status)
		if status == kFutureExpired {
			return predecessor
		}
		return last
	case kFutureExpired: // lastStatus == expired: 说明last还凑合着能用
		return last
	case kFutureRotted: // lastStatus == rotted: 说明last不能用了, 只能返回next
		return next
	}

	// lastStatus == empty: 也就是没有last, 所以只能返回next
	return next
}

func (my *cacheImpl) sendJob(job cacheJob) {
	// 如果Cache被Close()了, 通过closeChan确保不会因此被阻塞
	select {
	case my.jobChan <- job:
	case <-my.closeChan: // closeChan在NewCache()中已经初始化了
	}
}

//func (my *cacheImpl) Load(key interface{}, loader Loader) *Future {
//	assert(key != nil, "key is nil")
//	assert(loader != nil, "loader is nil")
//
//	var index, _ = cacheSharding.GetShardingIndex(key)
//	var futures = my.futures[index]
//
//	// 尝试获取缓存中的future
//	futures.RLock()
//	var last = futures.d[key]
//	futures.RUnlock()
//
//	var status = my.getFutureStatus(last)
//	if status == kFutureGood {
//		return last
//	}
//
//	// 以下代码需要考虑并发, 需要阻止重复加载
//	var next *Future = nil
//	futures.Lock()
//	{
//		last = futures.d[key]
//		status = my.getFutureStatus(last)
//		// 可能是empty, expired, rotted
//		if status != kFutureGood {
//			next = newFuture()
//			futures.d[key] = next
//			my.jobChan <- cacheJob{loader: loader, key: key, future: next}
//		}
//	}
//	futures.Unlock()
//
//	// 如果仅仅是expired, 但没到rotted状态, 则返回last, Get1()凑合用不用等IO
//	if status == kFutureGood || status == kFutureExpired {
//		return last
//	}
//
//	return next
//}

func (my *cacheImpl) removeRotted() {
	for _, futures := range my.futures {
		futures.Lock()
		for key, future := range futures.d {
			var status = my.getFutureStatus(future)
			if status == kFutureRotted {
				delete(futures.d, key)
			}
		}
		futures.Unlock()
	}
}

func (my *cacheImpl) getFutureStatus(future *Future) int {
	if future != nil {
		var updateTime = future.getUpdateTime()
		var past = time.Now().Sub(updateTime)

		var expire = my.args.normalExpire
		if future.err != nil {
			expire = my.args.errorExpire
		}

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
