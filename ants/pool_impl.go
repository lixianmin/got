package ants

import (
	"github.com/lixianmin/got/loom"
)

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type poolImpl struct {
	taskChan      chan Task
	taskChanInner chan Task
	closeChan     chan struct{}
}

func (my *poolImpl) Send(handler Handler, options ...TaskOption) Task {
	if handler == nil {
		panic("handler is nil")
	}

	var opts = createTaskOptions(options)
	var task = newTaskCallback(my, handler, opts)

	select {
	case my.taskChan <- task:
	case <-my.closeChan:
	}

	return task
}

func (my *poolImpl) sendTaskInner(task Task) {
	select {
	case my.taskChanInner <- task:
	case <-my.closeChan:
	}
}

func (my *poolImpl) goDispatchTask() {
	defer loom.DumpIfPanic()

	for {
		select {
		case task := <-my.taskChan:
			task.run()
		case <-my.closeChan:
			return
		}
	}
}

func (my *poolImpl) goDispatchTaskInner() {
	defer loom.DumpIfPanic()

	for {
		select {
		case task := <-my.taskChanInner:
			task.run()
		case <-my.closeChan:
			return
		}
	}
}
