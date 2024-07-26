package ants

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
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

	var task = pool.Send(func(ctx context.Context) (any, error) {
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
	var task = pool.Send(func(ctx context.Context) (any, error) {
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
	var pool = NewPool(WithSize(5))
	var startTime = time.Now()
	var task = pool.Send(func(ctx context.Context) (any, error) {
		time.Sleep(time.Second)
		return nil, nil
	}, WithTimeout(200*time.Millisecond), WithRetry(3))

	var _, err = task.Get2()

	var tasks = make([]Task, 0)
	for i := 0; i < 100; i++ {
		tasks = append(tasks, pool.Send(func(ctx context.Context) (any, error) {
			time.Sleep(10 * time.Millisecond)
			return nil, nil
		}, WithTimeout(5*time.Millisecond), WithRetry(3)))
	}

	for _, task := range tasks {
		task.Get1()
	}

	var endTime = time.Now()
	var past = endTime.Sub(startTime)
	if past > 2*time.Second || err != nil && !errors.Is(err, context.DeadlineExceeded) {
		t.Fail()
	}
}

func TestPool_ContextBuilder(t *testing.T) {
	var key = struct{}{}
	var counter = 0
	var pool = NewPool(WithSize(4), WithContextBuilder(func() context.Context {
		counter++
		return context.WithValue(context.Background(), key, counter)
	}))

	var wg sync.WaitGroup
	var size = 10
	wg.Add(size)

	for i := 0; i < size; i++ {
		_ = pool.Send(func(ctx context.Context) (any, error) {
			fmt.Printf("---> %v\n", ctx.Value(key))
			wg.Done()
			return nil, nil
		})
	}

	wg.Wait()
	runtime.GC()
}

func TestPool_DiscardOnBusy(t *testing.T) {
	var pool = NewPool()

	var wg sync.WaitGroup
	var size = 2
	wg.Add(size * 2)

	for i := 0; i < size*2; i++ {
		var task = pool.Send(func(ctx context.Context) (any, error) {
			time.Sleep(time.Second)
			wg.Done()
			return nil, nil
		})

		if discard, ok := task.(*taskDiscard); ok {
			wg.Done()
			var _, err = discard.Get2()
			fmt.Println(err)
		}
	}

	wg.Wait()
}
