package taskx

import (
	"github.com/lixianmin/got/loom"
	"github.com/lixianmin/got/std"
	"time"
)

/********************************************************************
created:    2020-09-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type delayedQueue struct {
	tasks chan *taskDelayed
}

func newDelayedQueue() *delayedQueue {
	var my = &delayedQueue{
		tasks: make(chan *taskDelayed, 128),
	}

	loom.Go(my.goLoop)
	return my
}

func (my *delayedQueue) goLoop(later loom.Later) {
	var ticker = later.NewTicker(1000 * time.Millisecond)
	var pq = std.NewPriorityQueue(32)

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

func (my *delayedQueue) PushTask(task *taskDelayed) {
	my.tasks <- task
}
