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

	var ctx, _ = context.WithTimeout(context.TODO(), time.Second)
	var task = pool.Send(ctx, func() (interface{}, error) {
		time.Sleep(time.Second / 2)
		return nil, nil
	})

	_, err := task.Get2()
	fmt.Println(err)
	runtime.GC()
}

func TestPool_Send2(t *testing.T) {
	var pool = NewPool(WithSize(8))

	var task = pool.Send(context.Background(), func() (interface{}, error) {
		time.Sleep(time.Second / 2)
		return nil, nil
	})

	task.Get1()
	task.Get1()
	task.Get1()
}
