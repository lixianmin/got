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
	var wc WaitClose
	var tasks = NewTaskQueue(TaskQueueArgs{
		Size:      8,
		CloseChan: wc.closeChan,
	})

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
			case <-wc.C():
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
	_ = wc.Close(nil)

	tasks.SendCallback(func(args interface{}) (result interface{}, err error) {
		println("oh oops")
		return nil, nil
	})
}

func TestTaskQueue_SendDelayed(t *testing.T) {
	var tasks = NewTaskQueue(TaskQueueArgs{})
	var closeChan = make(chan struct{})
	var delayedTime = 2 * time.Second

	tasks.SendDelayed(delayedTime, func(args interface{}) (i interface{}, e error) {
		fmt.Printf("--> args=%v, delayedTime=%s\n", args, delayedTime.String())
		time.Sleep(time.Second)
		close(closeChan)
		return nil, nil
	})

	go func() {
		for {
			select {
			case task := <-tasks.C:
				var err = task.Do(1)
				if err != nil {
					println(err)
				}
			case <-closeChan:
				break
			}
		}
	}()

	<-closeChan
}
