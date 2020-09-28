package std

import "container/heap"

/********************************************************************
created:    2020-09-28
author:     lixianmin

参考：https://github.com/gansidui/priority_queue

Copyright (C) - All Rights Reserved
*********************************************************************/

type PriorityQueueItem interface {
	Less(other interface{}) bool
}

type sorter []PriorityQueueItem

// Implement heap.Interface: Push, Pop, Len, Less, Swap
func (s *sorter) Push(x interface{}) {
	*s = append(*s, x.(PriorityQueueItem))
}

func (s *sorter) Pop() interface{} {
	n := len(*s)
	if n > 0 {
		x := (*s)[n-1]
		*s = (*s)[0 : n-1]
		return x
	}
	return nil
}

func (s *sorter) Len() int {
	return len(*s)
}

func (s *sorter) Less(i, j int) bool {
	return (*s)[i].Less((*s)[j])
}

func (s *sorter) Swap(i, j int) {
	(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
}

// Define priority queue struct
type PriorityQueue struct {
	s *sorter
}

func NewPriorityQueue() *PriorityQueue {
	var my = &PriorityQueue{
		s: new(sorter),
	}

	heap.Init(my.s)
	return my
}

func (my *PriorityQueue) Push(x PriorityQueueItem) {
	heap.Push(my.s, x)
}

func (my *PriorityQueue) Pop() PriorityQueueItem {
	return heap.Pop(my.s).(PriorityQueueItem)
}

func (my *PriorityQueue) Top() PriorityQueueItem {
	if len(*my.s) > 0 {
		return (*my.s)[0].(PriorityQueueItem)
	}

	return nil
}

func (my *PriorityQueue) Fix(x PriorityQueueItem, i int) {
	(*my.s)[i] = x
	heap.Fix(my.s, i)
}

func (my *PriorityQueue) Remove(i int) PriorityQueueItem {
	return heap.Remove(my.s, i).(PriorityQueueItem)
}

func (my *PriorityQueue) Len() int {
	return my.s.Len()
}
