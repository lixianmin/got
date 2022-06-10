package ants

import "runtime"

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type options struct {
	size int
}

type Option func(*options)

func createOptions(optionList []Option) options {
	var opts = options{
		size: runtime.NumCPU(),
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	return opts
}

func WithSize(size int) Option {
	return func(opts *options) {
		if size > 0 {
			opts.size = size
		}
	}
}
