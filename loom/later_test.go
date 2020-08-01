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
	wc *WaitClose
}

func (my *TestGo) Init() {
	my.wc = NewWaitClose()
	Go(my.goLoop)

	<-my.wc.CloseChan
}

func (my *TestGo) goLoop(later *Later) {
	var ticker = later.NewTicker(time.Second)
	var counter = 0
	for {
		select {
		case <-ticker.C:
			counter += 1
			println("goLoop")
			if counter == 3 {
				_ = my.wc.Close()
				return
			}
		}
	}
}

func TestLater_NewTicker(t *testing.T) {
	var test = &TestGo{}
	test.Init()
}
