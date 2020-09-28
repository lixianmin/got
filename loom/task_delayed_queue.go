package loom

import (
	"github.com/lixianmin/got/std"
	"time"
)

/********************************************************************
created:    2020-09-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskDelayedQueue struct {
	tasks chan *taskDelayed
}

func newTaskDelayedQueue() *taskDelayedQueue {
	var my = &taskDelayedQueue{
		tasks: make(chan *taskDelayed, 128),
	}

	Go(my.goLoop)
	return my
}

func (my *taskDelayedQueue) goLoop(later Later) {
	var ticker = later.NewTicker(1000 * time.Millisecond)
	var pq = std.NewPriorityQueue()

	for {
		select {
		case <-ticker.C:
			var now = time.Now()
			var timestamp = now.UnixNano()
			for pq.Len() > 0 {
				var task = pq.Top().(*taskDelayed)
				if task.triggerTime > timestamp {
					break
				}

				pq.Pop()
				_ = task.Do(nil)
			}
		case task := <-my.tasks:
			pq.Push(task)
		}
	}
}

func (my *taskDelayedQueue) PushTask(task *taskDelayed) {
	my.tasks <- task
}
