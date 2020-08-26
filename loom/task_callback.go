package loom

import "sync"

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskCallback struct {
	wg      sync.WaitGroup
	handler func() error
}

func (task *taskCallback) Do() error {
	var err = task.handler()
	task.wg.Done()
	return err
}

func (task *taskCallback) Wait() {
	task.wg.Wait()
}
