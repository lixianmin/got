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
	var cache = NewCache(WithParallel(4))
	var start = time.Now()

	var f1 = cache.Load(1, 10*time.Second, func(key interface{}) interface{} {
		time.Sleep(2)
		return 1
	})

	var f2 = cache.Load(2, 10*time.Second, func(key interface{}) interface{} {
		time.Sleep(4)
		return 2
	})

	var f3 = cache.Load(3, 10*time.Second, func(key interface{}) interface{} {
		time.Sleep(6)
		return 3
	})

	var f4 = cache.Load(4, 20*time.Second, func(key interface{}) interface{} {
		time.Sleep(4)
		return 4
	})

	var f5 = cache.Load(4, 20*time.Second, func(key interface{}) interface{} {
		time.Sleep(4)
		return 5
	})

	fmt.Printf("cost time = %d\n", time.Now().Sub(start)/time.Millisecond)
	if f4 != f5 {
		t.Fatalf("f4=%d, f5=%d", f4.Get(), f5.Get())
	}

	if f5.Get() != 4 {
		t.Fatalf("f5 should be the same as f4")
	}

	fmt.Printf("f1=%d, f2=%d, f3=%d, f4=%d, f5=%d\n", f1.Get(), f2.Get(), f3.Get(), f4.Get(), f5.Get())
	time.Sleep(2 * time.Second)
}
