package loom

import (
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
	var ticker3 = later.NewTicker(time.Second*2)
	ticker2.Stop()
	ticker3.Stop()
	var ticker4 = later.NewTicker(time.Hour)
	print(ticker4)
	//var timer = later.NewTimer(2 * time.Second)
	var counter = 0
	for {
		select {
		case <-ticker.C:
			counter += 1
			println("goLoop")
			if counter == 10 {
				my.wc.Close(nil)
				return
			}

			//case <-timer.C:
			//	var t2 = later.NewTimer(2 * time.Second)
			//	println("timer")
			//	println(t2)
		}
	}
}

func TestLater_NewTicker(t *testing.T) {
	var test = &TestGo{}
	test.Init()
}
