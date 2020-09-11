package loom

import (
	"fmt"
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
	w := NewTimingWheel(10*time.Millisecond, 101)

	for i := 0; i < 5; i++ {
		go func() {
			for {
				select {
				case <-w.After(time.Second):
					fmt.Println(time.Now().Unix())
				}
			}
		}()
	}

	select {}
}
