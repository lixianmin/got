package loom

import (
	"fmt"
	"github.com/lixianmin/got/randx"
	"testing"
	"time"
)

/********************************************************************
created:    2020-08-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TestGo struct {
	wc WaitClose
}

func (my *TestGo) Init() {
	Go(my.goLoop)

	<-my.wc.C()
}

func (my *TestGo) goLoop(later Later) {
	var ticker = later.NewTicker(time.Second)
	var ticker2 = later.NewTicker(time.Second)
	ticker2.Stop()

	var timer = later.NewTimer(2 * time.Second)
	var timer2 = later.NewTimer(time.Second)
	timer2.Stop()

	var timer3 = later.NewTimer(5 * time.Second)

	var counter = 0
	for {
		select {
		case <-ticker.C:
			counter += 1
			println("ticker")
			if counter == 10 {
				my.wc.Close(nil)
				return
			}

		case <-timer.C:
			timer.Reset(randx.Duration(0, 3*time.Second))
			fmt.Printf("timer: %s\n", time.Now().Format(time.RFC3339Nano))
		case <-timer3.C:
			println("timer3")
		}
	}
}

func TestLater_NewTicker(t *testing.T) {
	var test = &TestGo{}
	test.Init()
}
