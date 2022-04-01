package cachex

import (
	"time"
)

/********************************************************************
created:    2021-08-06
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type arguments struct {
	parallel     int
	normalExpire time.Duration
	errorExpire  time.Duration
	jobChanSize  int
}

type Option func(*arguments)

func createArguments(options []Option) arguments {
	var args = arguments{
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

func WithParallel(num int) Option {
	return func(args *arguments) {
		if num > 0 {
			args.parallel = num
		}
	}
}

func WithExpire(normal time.Duration, error time.Duration) Option {
	return func(args *arguments) {
		assert(normal >= error, "assert failed: normal>=error")
		assert(error > 0, "assert failed: error>0")
		args.normalExpire = normal
		args.errorExpire = error
	}
}

func WithJobChanSize(size int) Option {
	return func(args *arguments) {
		assert(size > 0, "assert failed: size>0")
		args.jobChanSize = size
	}
}
