package loom

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

/********************************************************************
created:    2020-07-13
author:     lixianmin

仿java.util.concurrent.ConcurrentMap实现的Map类，主要目标为：
1. 提供更高的写并发度
2. 提供像ComputeIfAbsent()这样的延迟初始化方法

Copyright (C) - All Rights Reserved
*********************************************************************/

var shardCount = fetchShardCount()
var shardCountMinus1 = shardCount - 1

type shardItem struct {
	sync.RWMutex
	items map[interface{}]interface{}
}

type Map struct {
	m    sync.Mutex
	data []*shardItem
	size int64
}

// 如果已经存在了相同key的value，则覆盖找返回以前存在的那一个值；否则返回nil
func (my *Map) Put(key interface{}, value interface{}) interface{} {
	var shard = my.getShard(key)
	var last interface{}
	var has = false
	shard.Lock()
	{
		last, has = shard.items[key]
		if has {
			shard.items[key] = value
		} else {
			shard.items[key] = value
			atomic.AddInt64(&my.size, 1)
		}
	}
	shard.Unlock()
	return last
}

// 如果存在，则删除，并返回该值
func (my *Map) Remove(key interface{}) interface{} {
	var shard = my.getShard(key)
	var last interface{}
	var has = false
	shard.Lock()
	{
		last, has = shard.items[key]
		if has {
			delete(shard.items, key)
			atomic.AddInt64(&my.size, -1)
		}
	}
	shard.Unlock()
	return last
}

// 如果map中存在，则返回；否则返回nil
func (my *Map) Get(key interface{}) interface{} {
	var shard = my.getShard(key)
	var last, _ = my.getInner(shard, key)
	return last
}

func (my *Map) getInner(shard *shardItem, key interface{}) (interface{}, bool) {
	var last interface{}
	var has = false
	shard.RLock()
	{
		last, has = shard.items[key]
	}
	shard.RUnlock()
	return last, has
}

// 这其实是一种get命令：如果有，直接返回； 如果没有，就放进去，然后返回
func (my *Map) PutIfAbsent(key interface{}, value interface{}) interface{} {
	var shard = my.getShard(key)
	var last, has = my.getInner(shard, key)
	if has {
		return last
	}

	shard.Lock()
	{
		last, has = shard.items[key]
		if !has {
			shard.items[key] = value
			atomic.AddInt64(&my.size, 1)
		}
	}
	shard.Unlock()
	return last
}

// 如果原来存在，则返回原来的值；否则使用creator创建一个新值，放到到map中，则返回它
func (my *Map) ComputeIfAbsent(key interface{}, creator func(key interface{}) interface{}) interface{} {
	var shard = my.getShard(key)
	var last, has = my.getInner(shard, key)
	if has {
		return last
	}

	my.m.Lock()
	defer my.m.Unlock() // 用defer是因为不知道creator会不会panic

	// 加x锁后需要重新测试有没有数据
	// 如果creator=nil，也就不会重新生成了
	last, has = shard.items[key]
	if has || creator == nil {
		return last
	}

	// 如果没有，则创建一个放入到容器中
	var item = creator(key)
	shard.items[key] = item
	atomic.AddInt64(&my.size, 1)
	return item
}

func (my *Map) Size() int {
	return int(atomic.LoadInt64(&my.size))
}

func (my *Map) getShard(key interface{}) *shardItem {
	var index = getShardIndex(key)
	if my.data == nil {
		my.m.Lock()
		if my.data == nil {
			my.data = make([]*shardItem, shardCount)
			for i := 0; i < shardCount; i++ {
				var item = &shardItem{items: make(map[interface{}]interface{}, 4)}
				my.data[i] = item
			}
		}
		my.m.Unlock()
	}

	var shard = my.data[index]
	return shard
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func getShardIndex(key interface{}) int {
	var next int
	switch key := key.(type) {
	case int:
		next = key
	case int8:
		next = int(key)
	case int16:
		next = int(key)
	case int32:
		next = int(key)
	case int64:
		next = int(key)
	case uint8:
		next = int(key)
	case uint16:
		next = int(key)
	case uint32:
		next = int(key)
	case uint64:
		next = int(key)
	case string:
		next = int(fnv32(key))
	default:
		var message = fmt.Sprintf("Not supported key type, key= %v", key)
		panic(message)
	}

	var index = next & shardCountMinus1
	return index
}

// 由于getShardIndex()算法需要，这个shardCount必须是 2 的指数倍
func fetchShardCount() int {
	var numCpu = runtime.NumCPU() << 1
	var result = 2
	for result < numCpu {
		result <<= 1
	}

	return result
}
