package taskx

import "time"

/********************************************************************
created:    2020-09-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskDelayed struct {
	queue       *Queue  // 任务所在的TaskQueue
	handler     Handler // handler
	triggerTime int64   // 触发任务的时间戳
}

func newTaskDelayed(queue *Queue, delayed time.Duration, handler Handler) *taskDelayed {
	var task = &taskDelayed{
		queue:       queue,
		handler:     handler,
		triggerTime: time.Now().Add(delayed).UnixNano(),
	}

	return task
}

func (task *taskDelayed) Do(args interface{}) error {
	task.queue.SendCallback(task.handler)
	return nil
}

func (task *taskDelayed) Less(other interface{}) bool {
	return task.triggerTime < other.(*taskDelayed).triggerTime
}
