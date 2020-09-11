package loom

import (
	"fmt"
	"io"
	"time"
)

/********************************************************************
created:    2020-08-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Later interface {
	NewTicker(d time.Duration) *time.Ticker
	NewTimer(d time.Duration) *time.Timer
}

type laterImpl struct {
	stoppers []interface{}
}

func (later *laterImpl) NewTicker(d time.Duration) *time.Ticker {
	var ticker = time.NewTicker(d)
	later.stoppers = append(later.stoppers, ticker)
	return ticker
}

func (later *laterImpl) NewTimer(d time.Duration) *time.Timer {
	var timer = time.NewTimer(d)
	later.stoppers = append(later.stoppers, timer)
	return timer
}

func (later *laterImpl) stop() {
	for i := len(later.stoppers) - 1; i >= 0; i-- {
		var item = later.stoppers[i]
		switch item := item.(type) {
		case *time.Ticker:
			item.Stop()
		case *time.Timer:
			item.Stop()
		case io.Closer:
			_ = item.Close()
		default:
			fmt.Printf("unknown item=%+v", item)
		}
	}
}

func Go(handler func(later Later)) {
	go func() {
		defer DumpIfPanic()

		var later = &laterImpl{}
		defer later.stop()

		handler(later)
	}()
}
