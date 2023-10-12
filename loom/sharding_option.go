package loom

import (
	"runtime"
)

/********************************************************************
created:    2021-02-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 有些cpu可能是64核的，创建太多的sharding数可能没有太大的意义
var defaultShardingCount = min(convertPowerOfTwo(runtime.NumCPU()), 32)

type shardingArguments struct {
	shardingCount int
}

type ShardingOption func(*shardingArguments)

func createShardingArguments(options []ShardingOption) shardingArguments {
	var args = shardingArguments{
		shardingCount: defaultShardingCount,
	}

	for _, opt := range options {
		opt(&args)
	}

	return args
}

func WithSharingCount(count int) ShardingOption {
	return func(args *shardingArguments) {
		if count > 0 {
			args.shardingCount = convertPowerOfTwo(count)
		}
	}
}

// 由于getShardingIndex()算法需要，这个shardCount必须是 2 的指数倍
func convertPowerOfTwo(sharding int) int {
	var result = 1
	for result < sharding {
		result <<= 1
	}

	return result
}
