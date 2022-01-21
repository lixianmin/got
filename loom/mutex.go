package loom

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

/********************************************************************
created:    2020-10-21
author:     lixianmin

参考：
	Go 并发编程实战课 04｜ Mutex：骇客编程，如何拓展额外功能？
	https://time.geekbang.org/column/article/296793

Copyright (C) - All Rights Reserved
*********************************************************************/

// 复制Mutex定义的常量
const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

// Mutex 扩展一个Mutex结构
type Mutex struct {
	sync.Mutex
}

// TryLock 尝试获取锁
func (m *Mutex) TryLock() bool {
	// fast path ：如果能成功抢到锁则直接返回
	// 下面这一句就是直接对应着sync.Mutex中Lock()方法，一模一样翻译过来的
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// 如果处于唤醒、加锁或者饥饿状态，这次请求就不参与竞争了，返回false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	// 尝试在竞争的状态下请求锁
	next := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, next)
}

func (m *Mutex) Count() int {
	// 获取state字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v >> mutexWaiterShift //得到等待者的数值
	v = v + (v & mutexLocked) //再加上锁持有者的数量，0或者1
	return int(v)
}

// IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexLocked == mutexLocked
}

// IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexWoken == mutexWoken
}

// IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return state&mutexStarving == mutexStarving
}

// WithLock 使用这个Mutex执行action
func (m *Mutex) WithLock(action func()) {
	if action != nil {
		m.Lock()
		defer m.Unlock()
		action()
	}
}
