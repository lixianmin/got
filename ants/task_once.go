package ants

import (
	"context"
)

/********************************************************************
created:    2022-06-15
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskOnce struct {
	ctx     context.Context
	handler Handler

	result   interface{}
	err      error
	doneChan chan struct{}
}

func newTaskOnce(ctx context.Context, handler Handler) *taskOnce {
	var my = &taskOnce{
		ctx:      ctx,
		handler:  handler,
		doneChan: make(chan struct{}),
	}

	return my
}

func (my *taskOnce) Get1() interface{} {
	var result, _ = my.Get2()
	return result
}

func (my *taskOnce) Get2() (interface{}, error) {
	select { // 调用Get2()与调用run()的不是同一个goroutine, 利用这一点可以节约一个goroutine
	case <-my.ctx.Done():
		return nil, context.DeadlineExceeded
	case <-my.doneChan:
		return my.result, my.err
	}
}

func (my *taskOnce) run() {
	defer close(my.doneChan)
	my.result, my.err = my.handler(my.ctx)
}
