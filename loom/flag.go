package loom

import "sync/atomic"

/********************************************************************
created:    2020-01-31
author:     lixianmin

	1. 使用int64直接定义的方式，而不是使用struct，是为了支持直接赋值：

	type Player {
		Flag
	}

	var player = &Player{}
	player.Flag = 0x0001

	2. 线程安全

Copyright (C) - All Rights Reserved
*********************************************************************/

type Flag int64

func (my *Flag) AddFlag(flag int64) {
	var addr = (*int64)(my)

	for {
		var last = atomic.LoadInt64(addr)
		var next = last | flag

		if atomic.CompareAndSwapInt64(addr, last, next) {
			break
		}
	}
}

func (my *Flag) RemoveFlag(flag int64) {
	var addr = (*int64)(my)

	for {
		var last = atomic.LoadInt64(addr)
		var next = last & ^flag

		if atomic.CompareAndSwapInt64(addr, last, next) {
			break
		}
	}
}

func (my *Flag) HasFlag(flag int64) bool {
	var addr = (*int64)(my)
	return (atomic.LoadInt64(addr) & flag) != 0
}
