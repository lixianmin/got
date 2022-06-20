package ants

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/
const minGoroutines = 2 //  因为task_callback会调用pool.send(), 因此只有一个goroutine的话会导致死锁

type poolOptions struct {
	size int
}

type PoolOption func(*poolOptions)

func createPoolOptions(optionList []PoolOption) poolOptions {
	var opts = poolOptions{
		size: minGoroutines,
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	return opts
}

func WithSize(size int) PoolOption {
	return func(opts *poolOptions) {
		if size >= minGoroutines {
			opts.size = size
		}
	}
}
