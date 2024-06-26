package std

import "container/heap"

/********************************************************************
created:    2020-09-28
author:     lixianmin

参考：https://github.com/gansidui/priority_queue

Copyright (C) - All Rights Reserved
*********************************************************************/

type Comparable interface {
	Less(other any) bool
}

type sorter []Comparable

// Push Implement heap.Interface: Push, Pop, Len, Less, Swap
func (s *sorter) Push(x any) {
	*s = append(*s, x.(Comparable))
}

func (s *sorter) Pop() any {
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

func NewPriorityQueue(capacity int) *PriorityQueue {
	var my = &PriorityQueue{
		s: new(sorter),
	}

	*my.s = make([]Comparable, 0, capacity)
	heap.Init(my.s)
	return my
}

func (my *PriorityQueue) Push(x Comparable) {
	heap.Push(my.s, x)
}

func (my *PriorityQueue) Pop() Comparable {
	return heap.Pop(my.s).(Comparable)
}

func (my *PriorityQueue) Top() Comparable {
	if len(*my.s) > 0 {
		return (*my.s)[0].(Comparable)
	}

	return nil
}

func (my *PriorityQueue) Fix(x Comparable, i int) {
	(*my.s)[i] = x
	heap.Fix(my.s, i)
}

func (my *PriorityQueue) Remove(i int) Comparable {
	return heap.Remove(my.s, i).(Comparable)
}

func (my *PriorityQueue) Len() int {
	return len(*my.s)
}
