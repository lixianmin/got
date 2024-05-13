package cachex

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

/********************************************************************
created:    2021-08-06
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestCache_Load(t *testing.T) {
	var cache = NewCache(WithParallel(4), WithExpire(10*time.Second, time.Second))
	var start = time.Now()

	var f1 = cache.Load(1, func(key any) (any, error) {
		time.Sleep(2 * time.Second)
		return 1, nil
	})

	var f2 = cache.Load(2, func(key any) (any, error) {
		time.Sleep(4 * time.Second)
		return 2, nil
	})

	var f3 = cache.Load(3, func(key any) (any, error) {
		time.Sleep(6 * time.Second)
		return 3, nil
	})

	var f4 = cache.Load(4, func(key any) (any, error) {
		time.Sleep(4)
		return 4, nil
	})

	var f5 = cache.Load(4, func(key any) (any, error) {
		time.Sleep(4 * time.Second)
		return 5, nil
	})

	fmt.Printf("cost time = %s\n", time.Now().Sub(start).String())
	if f4 != f5 {
		t.Fatalf("f4=%d, f5=%d", f4.Get1(), f5.Get1())
	}

	if f5.Get1() != 4 {
		t.Fatalf("f5 should be the same as f4")
	}

	fmt.Printf("f1=%d, f2=%d, f3=%d, f4=%d, f5=%d\n", f1.Get1(), f2.Get1(), f3.Get1(), f4.Get1(), f5.Get1())
	fmt.Printf("cost time = %s\n", time.Now().Sub(start).String())

	time.Sleep(time.Minute)
	//time.Sleep(time.Minute)
}

func TestCache_LoadMultiTimes(t *testing.T) {
	var cache = NewCache(WithParallel(4), WithExpire(time.Second, time.Second))
	var start = time.Now()

	var f1 = cache.Load(1, func(key interface{}) (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println(1)
		return 1, nil
	})

	cache.Load(1, func(key interface{}) (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println(2)
		return 2, nil
	})

	time.Sleep(4 * time.Second)

	for i := 0; i < 10; i++ {
		go cache.Load(1, func(key interface{}) (interface{}, error) {
			time.Sleep(time.Second)
			fmt.Println(3)
			return 3, nil
		})
	}

	var f2 = cache.Load(1, func(key interface{}) (interface{}, error) {
		time.Sleep(time.Second)
		fmt.Println(4)
		return 4, nil
	})

	fmt.Printf("f1=%d, f2=%d \n", f1.Get1(), f2.Get1())
	fmt.Printf("cost time = %s\n", time.Now().Sub(start).String())
}

func TestCache_Close(t *testing.T) {
	var cache = NewCache(WithJobChanSize(1))

	// 期望这里不会卡死
	for i := 0; i < 10; i++ {
		cache.Load(i, func(key interface{}) (interface{}, error) {
			return i, nil
		})
	}
}

func BenchmarkCache_LoadMultiTimes(t *testing.B) {
	var cache = NewCache(WithParallel(4), WithExpire(time.Microsecond, time.Microsecond/10))
	const threadCount = 100
	const loopCount = 1000

	var wg = sync.WaitGroup{}
	wg.Add(threadCount * loopCount)

	var err = errors.New("test")

	for i := 0; i < threadCount; i++ {
		go func() {
			for j := 0; j < loopCount; j++ {
				var k = j
				var future = cache.Load(k, func(key interface{}) (interface{}, error) {
					time.Sleep(time.Millisecond / 100)
					//fmt.Println(k)
					if rand.Float32() < 0.5 {
						return key, err
					}

					return key, nil
				})

				future.Get1()
				wg.Done()
			}
		}()
	}

	wg.Wait()
}

func TestCacheFinalizer(t *testing.T) {
	var cache = NewCache(WithParallel(4), WithExpire(time.Microsecond, time.Microsecond/10))
	var result = cache.Load("123", func(key interface{}) (interface{}, error) {
		return 10, nil
	})

	result.Get1()
	cache = nil

	runtime.GC()
	time.Sleep(1 * time.Second)
}

func TestCache_Predecessor(t *testing.T) {
	var cache = NewCache(WithParallel(4), WithExpire(2*time.Millisecond, time.Millisecond))

	var f1 = cache.Load(1, func(key interface{}) (interface{}, error) {
		return 1, nil
	}).Get1()

	time.Sleep(3 * time.Millisecond)
	var f2 = cache.Load(1, func(key interface{}) (interface{}, error) {
		return 2, nil
	}).Get1()

	var f3 = cache.Load(1, func(key interface{}) (interface{}, error) {
		return 3, nil
	}).Get1()

	time.Sleep(2 * time.Millisecond)
	var f4 = cache.Load(1, func(key interface{}) (interface{}, error) {
		return 4, nil
	}).Get1()

	fmt.Printf("f1=%d, f2=%d, f3=%d, f4=%d \n", f1, f2, f3, f4)
	if f1 != 1 || f2 != 1 || f3 != 1 || f4 != 2 {
		t.Fatal()
	}
}

func TestCache_GetSet(t *testing.T) {
	var cache = NewCache(WithExpire(2*time.Millisecond, time.Millisecond))
	cache.Set(1, 1, nil)
	var f1 = cache.Get(1)

	time.Sleep(3 * time.Millisecond)
	var f2 = cache.Get(1)

	time.Sleep(2 * time.Millisecond)
	var f3 = cache.Get(1)

	fmt.Printf("f1=%v, f2=%v, f3=%v \n", f1, f2, f3)
	if f1 != 1 || f2 != 1 || f3 != nil {
		t.Fatal()
	}
}
