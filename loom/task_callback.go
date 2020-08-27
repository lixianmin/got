package loom

import "sync"

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	handler func(args interface{}) (interface{}, error)
	wg      sync.WaitGroup
	result  interface{}
	err     error
}

func (task *taskCallback) Do(args interface{}) error {
	task.result, task.err = task.handler(args)
	task.wg.Done()
	return task.err
}

func (task *taskCallback) Get() (interface{}, error) {
	task.wg.Wait()
	return task.result, task.err
}
