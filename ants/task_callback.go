package ants

import (
	"context"
	"sync"
)

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	taskOptions
	handler Handler

	result interface{}
	err    error
	wg     sync.WaitGroup
}

func newTaskCallback(handler Handler, opts taskOptions) *taskCallback {
	var my = &taskCallback{
		taskOptions: opts,
		handler:     handler,
	}

	my.wg.Add(1)
	return my
}

func (my *taskCallback) Get1() interface{} {
	my.wg.Wait()
	return my.result
}

func (my *taskCallback) Get2() (interface{}, error) {
	my.wg.Wait()
	return my.result, my.err
}

func (my *taskCallback) run() {
	defer my.wg.Done()

	for i := 0; i < my.retry; i++ {
		if my.runOnce() == nil {
			return
		}
	}
}

func (my *taskCallback) runOnce() error {
	var ctx, cancel = context.WithTimeout(context.Background(), my.timeout)
	defer cancel()

	var doneChan = make(chan struct{})
	go func() {
		defer close(doneChan)
		my.result, my.err = my.handler(ctx)
	}()

	select {
	case <-doneChan:
	case <-ctx.Done():
		my.result, my.err = nil, context.DeadlineExceeded
	}

	return my.err
}
