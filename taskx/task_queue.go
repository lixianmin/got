package taskx

import (
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

type ITask interface {
	Do(args interface{}) error
	Get1() interface{}
	Get2() (interface{}, error)
}

type TaskQueue struct {
	closeChan chan struct{}
	C         chan ITask
}

func NewTaskQueue(options ...Option) *TaskQueue {
	var opts = createOptions(options)
	var my = &TaskQueue{
		closeChan: opts.CloseChan,
		C:         make(chan ITask, opts.Size),
	}

	return my
}

func (my *TaskQueue) SendTask(task ITask) ITask {
	if task != nil {
		select {
		case <-my.closeChan:
		case my.C <- task:
		}
	}

	return task
}

func (my *TaskQueue) SendCallback(handler TaskHandler) ITask {
	if handler == nil {
		return taskEmpty{}
	}

	var task = &taskCallback{handler: handler}
	task.wg.Add(1)

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
