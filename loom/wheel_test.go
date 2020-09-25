package loom

import (
	"fmt"
	"github.com/lixianmin/got/timex"
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

func TestWheelAfter(t *testing.T) {
	t.Parallel()
	wheel := NewWheel(10*time.Millisecond, 1001)

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

func TestWheelTicker(t *testing.T) {
	t.Parallel()
	wheel := NewWheel(500*time.Millisecond, 20+1)

	var count = 1
	var wg sync.WaitGroup
	wg.Add(count)

	go func() {
		defer wg.Done()

		var ticker1 = wheel.NewTicker(1 * time.Second)
		var ticker5 = wheel.NewTicker(5 * time.Second)

		for {
			select {
			case <-ticker1.C():
				fmt.Printf("every 1s : %s \n", timex.FormatTime(time.Now()))
			case <-ticker5.C():
				fmt.Printf("at the 5 second : %s \n", timex.FormatTime(time.Now()))
				return
			}
		}
	}()

	wg.Wait()
	_ = wheel.Close()
	_ = wheel.Close()
	_ = wheel.Close()
	//select {}
}
