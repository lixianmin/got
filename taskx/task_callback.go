package taskx

import "sync"

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	handler   func(args interface{}) (interface{}, error)
	wg        sync.WaitGroup
	result    interface{}
	err       error
	isHandled bool
}

// Do Do()方法通常是业务代码调用的, 因此可能会被重复调用多次
func (task *taskCallback) Do(args interface{}) error {
	task.result, task.err = task.handler(args)

	// 1. 只有第一次调用handler时调用wg.Done(), 防止多次调用导致panic
	// 2. 假定Do()方法只可能在同一个goroutine中反复调用, 不会跨多个goroutine调用, 不需要将isHandled设置为atomic变量
	// 3. handler()方法调用多次可能导致的幂等性问题, 由业务代码自己处理
	// 4. 如果有业务代码调用Get2()方法, 并且多次调用Do()方法的情况下, 无法保证Get2()返回的数据是最新的
	if !task.isHandled {
		task.isHandled = true
		task.wg.Done()
	}

	return task.err
}

func (task *taskCallback) Get1() interface{} {
	task.wg.Wait()
	return task.result
}

func (task *taskCallback) Get2() (interface{}, error) {
	task.wg.Wait()
	return task.result, task.err
}
