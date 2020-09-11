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
	w := NewTimingWheel(200*time.Millisecond, 10)

	for {
		select {
		case <-w.After(time.Second):
			fmt.Println(time.Now().Unix())
		}
	}
}
