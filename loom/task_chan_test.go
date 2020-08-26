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

	type Fetus struct {
		counter int
	}

	var fetus = Fetus{}

	go func() {
		for {
			select {
			case task := <-tasks.C:
				fetus.counter += 1
				var err = task.Do(fetus)
				if err != nil {
					println(err)
				}
			case <-wc.C:
				break
			}
		}
	}()

	tasks.SendCallback(func(args interface{}) error {
		var fetus = args.(Fetus)
		println("hello", fetus.counter)
		return nil
	})

	tasks.SendCallback(nil).Wait()

	tasks.SendCallback(func(args interface{}) error {
		time.Sleep(500 * time.Millisecond)
		var fetus = args.(Fetus)
		println("world", fetus.counter)
		return nil
	}).Wait()

	wc.Close()

	tasks.SendCallback(func(args interface{}) error {
		println("oh oops")
		return nil
	})
}
