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

func TestWheelTicker(t *testing.T) {
	t.Parallel()
	wheel := NewWheel(500*time.Millisecond, 20+1)

	var count = 1
	var wg sync.WaitGroup
	wg.Add(count)

	go func() {
		defer wg.Done()

		var timer1 = wheel.NewTimer(1 * time.Second)
		var timer5 = wheel.NewTimer(5 * time.Second)
		var startTime = time.Now()

		for {
			select {
			case <-timer1.C:
				timer1.Reset()
				fmt.Printf("every 1s : %s \n", time.Now().Format(timex.Layout))
			case <-timer5.C:
				timer5.Reset()
				fmt.Printf("at the 5 second : %s \n", time.Since(startTime).String())
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

func TestWheel_AfterFunc(t *testing.T) {
	t.Parallel()
	wheel := NewWheel(500*time.Millisecond, 20+1)

	var c = make(chan struct{})
	wheel.AfterFunc(2*time.Second, func() {
		fmt.Println("after 2s")
		close(c)
	})

	<-c
}
