package loom

import (
	"sync"
	"sync/atomic"
)

/********************************************************************
created:    2018-10-09
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type WaitClose struct {
	CloseChan chan struct{}
	m         sync.Mutex
	done      int32
}

func NewWaitClose() *WaitClose {
	var wd = &WaitClose{
		CloseChan: make(chan struct{}),
	}

	return wd
}

func (wd *WaitClose) Close() error {
	if 0 == atomic.LoadInt32(&wd.done) {
		wd.m.Lock()
		if 0 == wd.done {
			atomic.StoreInt32(&wd.done, 1)
			close(wd.CloseChan)
		}
		wd.m.Unlock()
	}

	return nil
}

func (wd *WaitClose) IsClosed() bool {
	return atomic.LoadInt32(&wd.done) == 1
}
