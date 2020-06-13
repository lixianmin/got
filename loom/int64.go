package loom

import "sync/atomic"

/********************************************************************
created:    2018-12-03
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Int64 int64

func (iam *Int64) Store(value int64) {
	atomic.StoreInt64((*int64)(iam), value)
}

func (iam *Int64) Load() int64 {
	return atomic.LoadInt64((*int64)(iam))
}

func (iam *Int64) AddIf(delta int64, predicate func(old int64) bool) bool {
	var expect, update int64
	for {
		expect = atomic.LoadInt64((*int64)(iam))
		if !predicate(expect) {
			return false
		}

		update = expect + delta

		if atomic.CompareAndSwapInt64((*int64)(iam), expect, update) {
			return true
		}
	}
}
