package loom

import "sync"

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	wg      sync.WaitGroup
	handler func(args interface{}) error
}

func (task *taskCallback) Do(args interface{}) error {
	var err = task.handler(args)
	task.wg.Done()
	return err
}

func (task *taskCallback) Wait() {
	task.wg.Wait()
}
