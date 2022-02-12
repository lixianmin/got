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
}

type CacheOption func(*cacheArguments)

func createCacheArguments(options []CacheOption) cacheArguments {
	var args = cacheArguments{
		parallel:     1,
		normalExpire: time.Second,
		errorExpire:  time.Millisecond * 100,
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
