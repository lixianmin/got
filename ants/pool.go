package ants

import (
	"context"
	"github.com/lixianmin/got/loom"
)

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Pool struct {
	taskChan chan *taskCallback
	wc       loom.WaitClose
}

func New(size int) *Pool {
	if size < 0 {
		panic("size < 0")
	}

	var my = &Pool{
		taskChan: make(chan *taskCallback, 8),
	}

	for i := 0; i < size; i++ {
		go my.goDispatch()
	}

	return my
}

func (my *Pool) Send(ctx context.Context, handler Handler) Task {
	if ctx == nil {
		panic("ctx is nil")
	}

	if handler == nil {
		panic("handler is nil")
	}

	var task = newTaskCallback(ctx, handler)
	select {
	case my.taskChan <- task:
	case <-my.wc.C():
	}

	return task
}

func (my *Pool) goDispatch() {
	defer loom.DumpIfPanic()
	var closeChan = my.wc.C()

	for {
		select {
		case task := <-my.taskChan:
			task.result, task.err = task.handler()
			task.doneChan <- struct{}{}
		case <-closeChan:
			return
		}
	}
}
