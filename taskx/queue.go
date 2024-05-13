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

var globalDelayedQueue = newDelayedQueue()

type Handler func(args any) (any, error)

type Queue struct {
	closeChan chan struct{}
	C         chan Task
	errLogger std.Logger
}

func NewQueue(options ...Option) *Queue {
	var opts = createOptions(options)
	var my = &Queue{
		closeChan: opts.closeChan,
		C:         make(chan Task, opts.size),
		errLogger: opts.errLogger,
	}

	return my
}

func (my *Queue) SendTask(task Task) Task {
	if task != nil {
		my.checkQueueFull()

		select {
		case <-my.closeChan:
		case my.C <- task:
		}
	}

	return task
}

func (my *Queue) SendCallback(handler Handler) Task {
	if handler == nil {
		return taskEmpty{}
	}

	var task = &taskCallback{handler: handler}
	task.wg.Add(1)
	my.checkQueueFull()

	select {
	case <-my.closeChan:
	case my.C <- task:
	}

	return task
}

func (my *Queue) SendDelayed(delayed time.Duration, handler Handler) {
	if handler == nil {
		return
	}

	var task = newTaskDelayed(my, delayed, handler)
	globalDelayedQueue.PushTask(task)
}

func (my *Queue) checkQueueFull() {
	var length = len(my.C)
	if length == cap(my.C) {
		my.errLogger.Printf("taskQueue is full, length=%d\n", length)
	}
}
