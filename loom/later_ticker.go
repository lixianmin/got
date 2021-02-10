package loom

import (
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2021-02-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type LaterTicker struct {
	*time.Ticker
	isStopped int32
}

func (my *LaterTicker) Stop() {
	my.Ticker.Stop()
	atomic.StoreInt32(&my.isStopped, 1)
}

func (my *LaterTicker) IsStopped() bool {
	return atomic.LoadInt32(&my.isStopped) == 1
}
