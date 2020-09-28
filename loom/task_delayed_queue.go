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
	queue *std.PriorityQueue
}

func newTaskDelayedQueue() *taskDelayedQueue {
	var my = &taskDelayedQueue{
		queue: std.NewPriorityQueue(),
	}

	Go(my.goLoop)
	return my
}

func (my *taskDelayedQueue) goLoop(later Later) {
	var ticker = later.NewTicker(1000 * time.Millisecond)
	var pq = my.queue

	for {
		select {
		case <-ticker.C:
			var timestamp = time.Now().UnixNano()
			for pq.Len() > 0 {
				var task = pq.Top().(*taskDelayed)
				if task.triggerTime < timestamp {
					break
				}

				pq.Pop()
				_ = task.Do(nil)
			}
		}
	}
}

func (my *taskDelayedQueue) PushTask(task *taskDelayed) {
	my.queue.Push(task)
}
