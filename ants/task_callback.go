package ants

import (
	"context"
	"sync"
)

/********************************************************************
created:    2022-06-15
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	taskOptions
	pool    *poolImpl
	handler Handler

	result any
	err    error
	wg     sync.WaitGroup
}

func newTaskCallback(pool *poolImpl, handler Handler, opts taskOptions) *taskCallback {
	var my = &taskCallback{
		taskOptions: opts,
		pool:        pool,
		handler:     handler,
	}

	my.wg.Add(1)
	return my
}

func (my *taskCallback) Get1() any {
	var result, _ = my.Get2()
	return result
}

func (my *taskCallback) Get2() (any, error) {
	my.wg.Wait()
	return my.result, my.err
}

func (my *taskCallback) run(ctx context.Context) {
	defer my.wg.Done()

	for i := 0; i < my.retry; i++ {
		my.runTaskOnce(ctx)
		if my.err == nil { // my.err是否是context.DeadlineExceeded, 都应该retry
			return
		}
	}
}

func (my *taskCallback) runTaskOnce(ctx context.Context) {
	var ctx1, cancel = context.WithTimeout(ctx, my.timeout)
	defer cancel() // cancel()使得my.handler(ctx)有时机检测到已经超时了, 可以提前返回

	var doneChan = make(chan struct{})
	my.pool.sendInnerCallback(func() {
		defer close(doneChan)
		var result, err = my.handler(ctx1)

		select {
		case <-ctx1.Done(): // 代码走到这里的时候, 一定是超时了, 外面的runTaskOnce()主体逻辑一定执行完成了, 因此不设置my.result
		default:
			my.result, my.err = result, err
		}
	})

	select {
	case <-doneChan:
	case <-ctx1.Done():
		my.result, my.err = nil, context.DeadlineExceeded
	}
}
