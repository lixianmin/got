package ants

import (
	"context"
	"runtime"
)

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Handler func() (interface{}, error)

type Pool interface {
	Send(ctx context.Context, handler Handler) Task
}

type wrapper struct {
	*poolImpl
}

func NewPool(options ...Option) Pool {
	var opts = createOptions(options)

	var my = &wrapper{&poolImpl{
		taskChan:  make(chan *taskCallback, opts.size),
		closeChan: make(chan struct{}),
	}}

	for i := 0; i < opts.size; i++ {
		go my.goDispatch()
	}

	// 参考: https://zhuanlan.zhihu.com/p/76504936
	runtime.SetFinalizer(my, func(w *wrapper) {
		close(w.closeChan) // 这里必须使用w.closeChan, 而不能使用my.closeChan, 否则runtime.GC()执行不到这里
		//fmt.Println("finalized")
	})

	return my
}