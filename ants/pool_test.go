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

func TestPool_Send2(t *testing.T) {
	var pool = NewPool(WithSize(8))
	var counter = 0
	var task = pool.Send(func(ctx context.Context) (interface{}, error) {
		counter++
		fmt.Println(counter)

		time.Sleep(time.Second)
		return nil, nil
	}, WithTimeout(time.Second), WithRetry(3))

	task.Get1()
	task.Get1()
	task.Get1()
}
