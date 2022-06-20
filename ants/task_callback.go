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
		if result, err := my.runTaskOnce(); err == nil || err != context.DeadlineExceeded {
			my.result, my.err = result, err
			return
		}
	}
}

func (my *taskCallback) runTaskOnce() (interface{}, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), my.timeout)
	defer cancel()

	var task = newTaskOnce(ctx, my.handler)
	my.pool.send(task)

	return task.Get2()
}
