package loom

import (
	"sync/atomic"
	"unsafe"
)

/********************************************************************
created:    2019-06-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func AddIf64(addr *int64, delta int64, predicate func(old int64) bool) bool {
	if addr == nil {
		return false
	}

	var expect, update int64
	for {
		expect = atomic.LoadInt64(addr)
		if !predicate(expect) {
			return false
		}

		update = expect + delta

		if atomic.CompareAndSwapInt64(addr, expect, update) {
			return true
		}
	}
}

func LoadString(addr *unsafe.Pointer) *string {
	return (*string)(atomic.LoadPointer(addr))
}

func StoreString(addr *unsafe.Pointer, s *string) *string {
	atomic.StorePointer(addr, unsafe.Pointer(s))
	return s
}
