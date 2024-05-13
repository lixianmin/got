package taskx

import (
	"fmt"
	"github.com/lixianmin/got/std"
	"os"
)

/********************************************************************
created:    2021-06-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type options struct {
	size      int
	closeChan chan struct{}
	errLogger std.Logger
}

type Option func(*options)

func createOptions(optionList []Option) options {
	var opts = options{
		size: 8,
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	// 默认创建一个长度为0的closeChan
	if opts.closeChan == nil {
		opts.closeChan = make(chan struct{})
	}

	// 默认创建一个打印到stderr的errLogger
	if opts.errLogger == nil {
		opts.errLogger = std.LoggerFunc(func(format string, args ...any) {
			var message = fmt.Sprintf(format, args...)
			_, _ = os.Stderr.WriteString(message)
		})
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

func WithCloseChan(c chan struct{}) Option {
	return func(opts *options) {
		if c != nil {
			opts.closeChan = c
		}
	}
}

func WithErrorLogger(logger std.LoggerFunc) Option {
	return func(opts *options) {
		if logger != nil {
			opts.errLogger = logger
		}
	}
}
