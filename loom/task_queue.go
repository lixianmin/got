package loom

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ITask interface {
	Do() error
	Wait()
}

type TaskQueue struct {
	closeChan chan struct{}
	TaskChan  chan ITask
}

func NewTaskQueue(closeChan chan struct{}) *TaskQueue {
	if closeChan == nil {
		closeChan = make(chan struct{})
	}

	var my = &TaskQueue{
		closeChan: closeChan,
		TaskChan:  make(chan ITask, 8),
	}

	return my
}

func (my *TaskQueue) SendTask(task ITask) ITask {
	if task != nil {
		select {
		case <-my.closeChan:
		case my.TaskChan <- task:
		}
	}

	return task
}

func (my *TaskQueue) SendCallback(handler func() error) ITask {
	if handler == nil {
		return taskEmpty{}
	}

	var task = &taskCallback{handler: handler}
	task.wg.Add(1)

	select {
	case <-my.closeChan:
	case my.TaskChan <- task:
	}

	return task
}
