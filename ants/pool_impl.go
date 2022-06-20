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
	taskChan  chan Task
	closeChan chan struct{}
}

func (my *poolImpl) Send(handler Handler, options ...TaskOption) Task {
	if handler == nil {
		panic("handler is nil")
	}

	var opts = createTaskOptions(options)
	var task = newTaskCallback(my, handler, opts)
	my.send(task)
	return task
}

func (my *poolImpl) send(task Task) {
	select {
	case my.taskChan <- task:
	case <-my.closeChan:
	}
}

func (my *poolImpl) goDispatch() {
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
