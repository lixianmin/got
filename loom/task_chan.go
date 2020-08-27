package loom

/********************************************************************
created:    2020-08-25
author:     lixianmin

改名为TaskChan是为了跟以后可能出现的TaskQueue相区别：
1. TaskChan只负责接收任务，但并不负责消费任务
2. TaskQueue应该会自己启动一个独立的goroutine消费掉所有任务

Copyright (C) - All Rights Reserved
*********************************************************************/

type ITask interface {
	Do(args interface{}) error
	Get1() interface{}
	Get2() (interface{}, error)
}

type TaskChan struct {
	closeChan chan struct{}
	C         chan ITask
}

func NewTaskChan(closeChan chan struct{}) *TaskChan {
	if closeChan == nil {
		closeChan = make(chan struct{})
	}

	var my = &TaskChan{
		closeChan: closeChan,
		C:         make(chan ITask, 8),
	}

	return my
}

func (my *TaskChan) SendTask(task ITask) ITask {
	if task != nil {
		select {
		case <-my.closeChan:
		case my.C <- task:
		}
	}

	return task
}

func (my *TaskChan) SendCallback(handler func(args interface{}) (result interface{}, err error)) ITask {
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
