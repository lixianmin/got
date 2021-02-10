package loom

import (
	"time"
)

/********************************************************************
created:    2020-08-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Later interface {
	NewTicker(d time.Duration) *LaterTicker
}

type stopper interface {
	Stop()
	IsStopped() bool
}

type laterImpl struct {
	stoppers []stopper
}

func (later *laterImpl) NewTicker(d time.Duration) *LaterTicker {
	var ticker = &LaterTicker{
		Ticker:    time.NewTicker(d),
		isStopped: 0,
	}

	later.stoppers = append(tailorStopped(later.stoppers), ticker)
	return ticker
}

func tailorStopped(data []stopper) []stopper {
	var size = len(data)
	if size > 0 {
		var i, j = 0, len(data) - 1
		for {
			for i <= j && !data[i].IsStopped() {
				i++
			}

			for i < j && data[j].IsStopped() {
				j--
			}

			if i < j {
				data[i], data[j] = data[j], data[i]
			} else {
				break
			}
		}

		data = data[:i]
	}

	return data
}

func (later *laterImpl) stop() {
	for i := len(later.stoppers) - 1; i >= 0; i-- {
		var item = later.stoppers[i]
		item.Stop()
	}
}

// loom.Go()，会启动一个协程，并在defer中调用stop()
func Go(handler func(later Later)) {
	go func() {
		defer DumpIfPanic()

		var later = &laterImpl{}
		defer later.stop()

		handler(later)
	}()
}
