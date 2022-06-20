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

	result interface{}
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

func (my *taskCallback) Get1() interface{} {
	var result, _ = my.Get2()
	return result
}

func (my *taskCallback) Get2() (interface{}, error) {
	my.wg.Wait()
	return my.result, my.err
}

func (my *taskCallback) run() {
	defer my.wg.Done()

	for i := 0; i < my.retry; i++ {
		my.result, my.err = my.runTaskOnce()
		if my.err == nil { // my.err是否是context.DeadlineExceeded, 都需要重试
			return
		}
	}
}

func (my *taskCallback) runTaskOnce() (interface{}, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), my.timeout)
	defer cancel()

	// 这个run()是在goDispatch()的goroutine中, 自己给自己发task, 很容易死锁
	var task = newTaskOnce(ctx, my.handler)
	my.pool.sendTaskInner(task)

	return task.Get2()
}
