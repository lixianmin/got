package loom

import (
	"fmt"
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

	tasks.SendCallback(func(args interface{}) (result interface{}, err error) {
		var fetus = args.(Fetus)
		println("hello", fetus.counter)
		return nil, nil
	})

	tasks.SendCallback(nil).Get1()

	result, _ := tasks.SendCallback(func(args interface{}) (interface{}, error) {
		time.Sleep(500 * time.Millisecond)
		var fetus = args.(Fetus)
		result := fmt.Sprintf("world %d", fetus.counter)
		return result, nil
	}).Get2()

	println(result.(string))
	wc.Close()

	tasks.SendCallback(func(args interface{}) (result interface{}, err error) {
		println("oh oops")
		return nil, nil
	})
}
