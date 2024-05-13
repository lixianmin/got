package std

import (
	"fmt"
	"testing"
)

/********************************************************************
created:    2020-09-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Item struct {
	priority int
}

func (my *Item) Less(other any) bool {
	return my.priority < other.(*Item).priority
}

func TestTaskDelayedQueue(t *testing.T) {
	var items = []*Item{
		{priority: 1},
		{priority: 4},
		{priority: 9},
		{priority: 5},
		{priority: 3},
		{priority: 15},
		{priority: 2},
		{priority: 10},
		{priority: 7},
		{priority: 3},
	}

	var pq = NewPriorityQueue(8)
	for _, item := range items {
		pq.Push(item)
	}

	for pq.Len() > 0 {
		var item = pq.Pop().(*Item)
		fmt.Printf("%d \n", item.priority)
	}
}
