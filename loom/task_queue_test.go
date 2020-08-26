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
	var tq = NewTaskQueue(wc.C)

	go func() {
		for {
			select {
			case task := <-tq.TaskChan:
				var err = task.Do()
				if err != nil {
					println(err)
				}
			case <-wc.C:
				break
			}
		}
	}()

	tq.SendCallback(func() error {
		println("hello")
		return nil
	})

	tq.SendCallback(nil).Wait()

	tq.SendCallback(func() error {
		time.Sleep(500 * time.Millisecond)
		println("world")
		return nil
	}).Wait()

	wc.Close()

	tq.SendCallback(func() error {
		println("oh oops")
		return nil
	})
}
