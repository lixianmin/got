package ants

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

/********************************************************************
created:    2022-06-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestPool_Send(t *testing.T) {
	var pool = NewPool(WithSize(8))

	var task = pool.Send(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second / 2)
		return nil, nil
	}, WithTimeout(time.Second))

	_, err := task.Get2()
	fmt.Println(err)
	runtime.GC()
}

func TestPool_GetMultiTimes(t *testing.T) {
	const retry = 3
	var pool = NewPool(WithSize(8))
	var counter = 0
	var task = pool.Send(func(ctx context.Context) (interface{}, error) {
		counter++
		fmt.Println(counter)

		time.Sleep(time.Second)
		return nil, nil
	}, WithTimeout(time.Second), WithRetry(retry))

	task.Get1()

	if counter == retry {
		t.Fail()
	}

	task.Get1()
	task.Get1()
}

func TestPool_HandleTooLongTime(t *testing.T) {
	var pool = NewPool()
	var startTime = time.Now()
	var task = pool.Send(func(ctx context.Context) (interface{}, error) {
		time.Sleep(time.Second)
		return nil, nil
	}, WithTimeout(200*time.Millisecond))

	var _, err = task.Get2()
	var endTime = time.Now()
	var past = endTime.Sub(startTime)
	if past > 500*time.Millisecond || err != context.DeadlineExceeded {
		t.Fail()
	}
}
