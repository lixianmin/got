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

type poolImpl struct {
	taskChan  chan *taskCallback
	closeChan chan struct{}
}

func (my *poolImpl) Send(ctx context.Context, handler Handler) Task {
	if ctx == nil {
		panic("ctx is nil")
	}

	if handler == nil {
		panic("handler is nil")
	}

	var task = newTaskCallback(ctx, handler)
	select {
	case my.taskChan <- task:
	case <-my.closeChan:
	}

	return task
}

func (my *poolImpl) goDispatch() {
	defer loom.DumpIfPanic()

	for {
		select {
		case task := <-my.taskChan:
			task.result, task.err = task.handler()
			close(task.doneChan)
		case <-my.closeChan:
			return
		}
	}
}
