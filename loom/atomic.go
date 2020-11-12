package loom

import (
	"sync/atomic"
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

// 后续基本数据类型的原子操作建议使用 https://github.com/uber-go/atomic ，不再提供支持
//func LoadString(addr **string) string {
//	// 最内层的unsafe.Pointer(addr)对第三方指针的包装，相当时强制类型转换，是固定写法
//	var p = (*string)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(addr))))
//	if p != nil {
//		return *p
//	}
//
//	return ""
//}
//
//func StoreString(addr **string, s string) string {
//	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(addr)), unsafe.Pointer(&s))
//	return s
//}
//
//func LoadBool(addr *int32) bool {
//	var v = atomic.LoadInt32(addr)
//	return v != 0
//}
//
//func StoreBool(addr *int32, b bool) bool {
//	var v int32 = 0
//	if b {
//		v = 1
//	}
//
//	atomic.StoreInt32(addr, v)
//	return b
//}
