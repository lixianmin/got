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
	taskChan          chan Task
	innerCallbackChan chan func()
	closeChan         chan struct{}
}

func (my *poolImpl) Send(handler Handler, opts ...TaskOption) Task {
	if handler == nil {
		panic("handler is nil")
	}

	var options = createTaskOptions(opts)
	if options.discardOnBusy && len(my.taskChan) == cap(my.taskChan) {
		var onError = options.onError
		if onError != nil {
			onError(errDiscard)
		}

		return newTaskDiscard()
	}

	var task = newTaskCallback(my, handler, options)
	select {
	case my.taskChan <- task:
	case <-my.closeChan:
	}

	return task
}

func (my *poolImpl) sendInnerCallback(callback func()) {
	select {
	case my.innerCallbackChan <- callback:
	case <-my.closeChan:
	}
}

func (my *poolImpl) goDispatchTask(ctx context.Context) {
	defer loom.DumpIfPanic()

	for {
		select {
		case task := <-my.taskChan:
			task.run(ctx)
		case <-my.closeChan:
			return
		}
	}
}

func (my *poolImpl) goDispatchInnerCallback() {
	defer loom.DumpIfPanic()

	for {
		select {
		case callback := <-my.innerCallbackChan:
			callback()
		case <-my.closeChan:
			return
		}
	}
}
