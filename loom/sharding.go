package loom

import (
	"fmt"
)

/********************************************************************
created:    2021-02-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Sharding struct {
	shardingCount int
}

func NewSharding(opts ...ShardingOption) *Sharding {
	var args = createShardingArguments(opts)
	var my = &Sharding{
		shardingCount: args.shardingCount,
	}

	return my
}

func (my *Sharding) GetShardingCount() int {
	return my.shardingCount
}

func (my *Sharding) GetShardingIndex(key any) (index int, normalizedKey any) {
	var shardingCount = my.shardingCount

	var next int64
	switch key := key.(type) {
	case int:
		next = int64(key)
	case int8:
		next = int64(key)
	case int16:
		next = int64(key)
	case int32:
		next = int64(key)
	case int64:
		next = key
	case uint8:
		next = int64(key)
	case uint16:
		next = int64(key)
	case uint32:
		next = int64(key)
	case uint64:
		next = int64(key)
	case string:
		next = int64(fnv32(key))
		index = int(next) & (shardingCount - 1)
		normalizedKey = key
		return
	default:
		var message = fmt.Sprintf("Not supported key type, key= %v", key)
		panic(message)
	}

	index = int(next) & (shardingCount - 1)
	normalizedKey = next
	return
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
