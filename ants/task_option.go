package ants

import (
	"github.com/lixianmin/got/timex"
	"time"
)

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskOptions struct {
	timeout time.Duration
	retry   int
}

type TaskOption func(*taskOptions)

func createTaskOptions(optionList []TaskOption) taskOptions {
	var opts = taskOptions{
		timeout: 365 * timex.Day, // 默认给一个∞
		retry:   1,
	}

	for _, opt := range optionList {
		opt(&opts)
	}

	return opts
}

func WithTimeout(timeout time.Duration) TaskOption {
	return func(opts *taskOptions) {
		if timeout > 0 {
			opts.timeout = timeout
		}
	}
}

func WithRetry(count int) TaskOption {
	return func(opts *taskOptions) {
		if count > 0 {
			opts.retry = count
		}
	}
}
