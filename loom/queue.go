package loom

/********************************************************************
created:    2021-06-02
author:     lixianmin

使用 Go 实现 lock-free 的队列:
https://colobu.com/2020/08/14/lock-free-queue-in-go/

github地址：
https://github.com/smallnest/queue

Copyright (C) - All Rights Reserved
*********************************************************************/

import (
	"sync/atomic"
	"unsafe"
)

// Queue is a lock-free unbounded queue.
type Queue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

type node struct {
	value any
	next  unsafe.Pointer
}

// NewQueue returns an empty queue.
func NewQueue() *Queue {
	n := unsafe.Pointer(&node{})
	return &Queue{head: n, tail: n}
}

// Push puts the given value v at the tail of the queue.
func (q *Queue) Push(v any) {
	n := &node{value: v}
	for {
		tail := queueLoad(&q.tail)
		next := queueLoad(&tail.next)
		if tail == queueLoad(&q.tail) { // are tail and next consistent?
			if next == nil {
				if queueCas(&tail.next, next, n) {
					queueCas(&q.tail, tail, n) // Push is done.  try to swing tail to the inserted node
					return
				}
			} else { // tail was not pointing to the last node
				// try to swing Tail to the next node
				queueCas(&q.tail, tail, next)
			}
		}
	}
}

// Pop removes and returns the value at the head of the queue.
// It returns nil if the queue is empty.
func (q *Queue) Pop() any {
	for {
		head := queueLoad(&q.head)
		tail := queueLoad(&q.tail)
		next := queueLoad(&head.next)
		if head == queueLoad(&q.head) { // are head, tail, and next consistent?
			if head == tail { // is queue empty or tail falling behind?
				if next == nil { // is queue empty?
					return nil
				}

				// tail is falling behind.  try to advance it
				queueCas(&q.tail, tail, next)
			} else {
				// read value before CAS otherwise another dequeue might free the next node
				v := next.value
				if queueCas(&q.head, head, next) {
					return v // Pop is done.  return
				}
			}
		}
	}
}

func queueLoad(p *unsafe.Pointer) (n *node) {
	return (*node)(atomic.LoadPointer(p))
}

func queueCas(p *unsafe.Pointer, old, new *node) (ok bool) {
	return atomic.CompareAndSwapPointer(p, unsafe.Pointer(old), unsafe.Pointer(new))
}
