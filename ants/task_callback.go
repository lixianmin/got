package ants

import "context"

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	ctx     context.Context
	handler Handler

	result   interface{}
	err      error
	doneChan chan struct{}
}

func newTaskCallback(ctx context.Context, handler Handler) *taskCallback {
	var my = &taskCallback{
		ctx:      ctx,
		handler:  handler,
		doneChan: make(chan struct{}),
	}

	return my
}

func (my *taskCallback) Get1() interface{} {
	var result, _ = my.Get2()
	return result
}

func (my *taskCallback) Get2() (interface{}, error) {
	select {
	case <-my.ctx.Done():
		return nil, context.DeadlineExceeded
	case <-my.doneChan:
		return my.result, my.err
	}
}
