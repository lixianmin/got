package cachex

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

/********************************************************************
created:    2021-08-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Future struct {
	value       any
	err         error
	updateTime  unsafe.Pointer
	predecessor unsafe.Pointer
	wg          sync.WaitGroup
}

func newFuture(predecessor *Future) *Future {
	var item = &Future{}
	atomic.StorePointer(&item.updateTime, unsafe.Pointer(&time.Time{}))
	atomic.StorePointer(&item.predecessor, unsafe.Pointer(predecessor))

	item.wg.Add(1)
	return item
}

func (my *Future) Get1() any {
	my.wg.Wait()
	return my.value
}

func (my *Future) Get2() (any, error) {
	my.wg.Wait()
	return my.value, my.err
}

// 这个方法只会被调用一次, 这个特别重要, 因为wg.Done()如果调用2次会panic
func (my *Future) setValue(value any, err error) {
	my.value = value
	my.err = err

	var now = time.Now()
	atomic.StorePointer(&my.updateTime, unsafe.Pointer(&now))
	atomic.StorePointer(&my.predecessor, nil)

	my.wg.Done()
}

func (my *Future) getUpdateTime() time.Time {
	var p = (*time.Time)(atomic.LoadPointer(&my.updateTime))
	if p != nil {
		return *p
	}

	return time.Time{}
}

func (my *Future) getPredecessor() *Future {
	var p = (*Future)(atomic.LoadPointer(&my.predecessor))
	return p
}
