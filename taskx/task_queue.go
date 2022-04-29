package taskx

import (
	"github.com/lixianmin/got/std"
	"time"
)

/********************************************************************
created:    2020-08-25
author:     lixianmin

改名为TaskChan是为了跟以后可能出现的TaskQueue相区别：
1. TaskChan只负责接收任务，但并不负责消费任务
2. TaskQueue应该会自己启动一个独立的goroutine消费掉所有任务

Copyright (C) - All Rights Reserved
*********************************************************************/

var globalTaskDelayedQueue = newTaskDelayedQueue()

type TaskHandler func(args interface{}) (interface{}, error)

type TaskQueue struct {
	closeChan chan struct{}
	C         chan Task
	errLogger std.Logger
}

func NewTaskQueue(options ...Option) *TaskQueue {
	var opts = createOptions(options)
	var my = &TaskQueue{
		closeChan: opts.closeChan,
		C:         make(chan Task, opts.size),
		errLogger: opts.errLogger,
	}

	return my
}

func (my *TaskQueue) SendTask(task Task) Task {
	if task != nil {
		my.checkTaskQueueFull()

		select {
		case <-my.closeChan:
		case my.C <- task:
		}
	}

	return task
}

func (my *TaskQueue) SendCallback(handler TaskHandler) Task {
	if handler == nil {
		return taskEmpty{}
	}

	var task = &taskCallback{handler: handler}
	task.wg.Add(1)
	my.checkTaskQueueFull()

	select {
	case <-my.closeChan:
	case my.C <- task:
	}

	return task
}

func (my *TaskQueue) SendDelayed(delayed time.Duration, handler TaskHandler) {
	if handler == nil {
		return
	}

	var task = newTaskDelayed(my, delayed, handler)
	globalTaskDelayedQueue.PushTask(task)
}

func (my *TaskQueue) checkTaskQueueFull() {
	var length = len(my.C)
	if length == cap(my.C) {
		my.errLogger.Printf("taskQueue is full, length=%d", length)
	}
}
