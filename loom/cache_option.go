package loom

import "time"

/********************************************************************
created:    2021-08-06
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type cacheArguments struct {
	parallel     int
	normalExpire time.Duration
	errorExpire  time.Duration
	jobChanSize  int
}

type CacheOption func(*cacheArguments)

func createCacheArguments(options []CacheOption) cacheArguments {
	var args = cacheArguments{
		parallel:     1,
		normalExpire: time.Second,
		errorExpire:  time.Millisecond * 100,
		jobChanSize:  128, // 加大这个chan的长度, 有助于减小第一次checkLoad()时的执行时间
	}

	for _, opt := range options {
		opt(&args)
	}

	return args
}

func WithParallel(num int) CacheOption {
	return func(args *cacheArguments) {
		if num > 0 {
			args.parallel = num
		}
	}
}

func WithExpire(normal time.Duration, error time.Duration) CacheOption {
	return func(args *cacheArguments) {
		assert(normal >= error, "assert failed: normal>=error")
		assert(error > 0, "assert failed: error>0")
		args.normalExpire = normal
		args.errorExpire = error
	}
}

func WithJobChanSize(size int) CacheOption {
	return func(args *cacheArguments) {
		assert(size > 0, "assert failed: size>0")
		args.jobChanSize = size
	}
}
