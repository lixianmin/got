package loom

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

/********************************************************************
created:    2020-09-11
author:     lixianmin

参考：https://blog.csdn.net/siddontang/article/details/18370541
https://github.com/siddontang/go/tree/master/timingwheel

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestTimingWheel(t *testing.T) {
	t.Parallel()
	wheel := NewTimingWheel(10*time.Millisecond, 1001)

	var count = 5
	var wg sync.WaitGroup
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			var exit = wheel.After(10 * time.Second)
			for {
				select {
				case <-wheel.After(1 * time.Second):
					fmt.Println(time.Now().Unix())
				case <-exit:
					return
				}
			}
		}()
	}

	wg.Wait()
	_ = wheel.Close()
	_ = wheel.Close()
	_ = wheel.Close()
	//select {}
}
