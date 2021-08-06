package loom

import (
	"fmt"
	"testing"
	"time"
)

/********************************************************************
created:    2021-08-06
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestCache_Load(t *testing.T) {
	var cache = NewCache(WithParallel(4), WithExpire(10*time.Second))
	var start = time.Now()

	var f1 = cache.Load(1, func(key interface{}) interface{} {
		time.Sleep(2 * time.Second)
		return 1
	})

	var f2 = cache.Load(2, func(key interface{}) interface{} {
		time.Sleep(4 * time.Second)
		return 2
	})

	var f3 = cache.Load(3, func(key interface{}) interface{} {
		time.Sleep(6 * time.Second)
		return 3
	})

	var f4 = cache.Load(4, func(key interface{}) interface{} {
		time.Sleep(4)
		return 4
	})

	var f5 = cache.Load(4, func(key interface{}) interface{} {
		time.Sleep(4 * time.Second)
		return 5
	})

	fmt.Printf("cost time = %s\n", time.Now().Sub(start).String())
	if f4 != f5 {
		t.Fatalf("f4=%d, f5=%d", f4.Get(), f5.Get())
	}

	if f5.Get() != 4 {
		t.Fatalf("f5 should be the same as f4")
	}

	fmt.Printf("f1=%d, f2=%d, f3=%d, f4=%d, f5=%d\n", f1.Get(), f2.Get(), f3.Get(), f4.Get(), f5.Get())
	fmt.Printf("cost time = %s\n", time.Now().Sub(start).String())

	time.Sleep(time.Minute)
	time.Sleep(time.Minute)
}
