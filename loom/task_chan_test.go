package loom

import (
	"testing"
	"time"
)

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestNewTaskQueue(t *testing.T) {

	var wc = NewWaitClose()
	var tasks = NewTaskChan(wc.C)

	go func() {
		for {
			select {
			case task := <-tasks.C:
				var err = task.Do()
				if err != nil {
					println(err)
				}
			case <-wc.C:
				break
			}
		}
	}()

	tasks.SendCallback(func() error {
		println("hello")
		return nil
	})

	tasks.SendCallback(nil).Wait()

	tasks.SendCallback(func() error {
		time.Sleep(500 * time.Millisecond)
		println("world")
		return nil
	}).Wait()

	wc.Close()

	tasks.SendCallback(func() error {
		println("oh oops")
		return nil
	})
}
