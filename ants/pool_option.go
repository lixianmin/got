package ants

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type poolOptions struct {
	size int
}

type PoolOption func(*poolOptions)

func createPoolOptions(optionList []PoolOption) poolOptions {
	var opts = poolOptions{
		size: 1,
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	return opts
}

func WithSize(size int) PoolOption {
	return func(opts *poolOptions) {
		if size > 0 {
			opts.size = size
		}
	}
}
