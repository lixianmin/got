package loom

import "time"

/********************************************************************
created:    2021-08-06
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type cacheArguments struct {
	parallel   int
	expire     time.Duration
	gcInterval time.Duration
}

type CacheOption func(*cacheArguments)

func createCacheArguments(options []CacheOption) cacheArguments {
	var args = cacheArguments{
		parallel:   1,
		expire:     time.Second,
		gcInterval: time.Minute,
	}

	for _, opt := range options {
		opt(&args)
	}

	return args
}

func WithParallel(count int) CacheOption {
	return func(args *cacheArguments) {
		if count > 0 {
			args.parallel = count
		}
	}
}

func WithExpire(expire time.Duration) CacheOption {
	return func(args *cacheArguments) {
		if expire > 0 {
			args.expire = expire
		}
	}
}

func WithGCInterval(interval time.Duration) CacheOption {
	return func(args *cacheArguments) {
		if interval > 0 {
			args.gcInterval = interval
		}
	}
}
