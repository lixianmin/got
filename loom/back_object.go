package loom

import (
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

/********************************************************************
created:    2019-01-31
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type BackObject struct {
	lock         sync.Mutex
	done         int32
	loadInterval time.Duration
	loader       func() (interface{}, error)
	data         atomic.Value
}

// NewBackObject 将loader()放到构造方法中，而不是将原来那样放到Get()方法中，是因为需要把loader()的closure固定下来。在Get()中有一个风险是
// 可能会无意中在closure中使用了变化的参数。因为loader()是周期性调用的，因此不应该使用可变的参数。
// 另外，每次调用Get()时创建新的closure也是一个额外的开销
func NewBackObject(loadInterval time.Duration, loader func() (interface{}, error)) *BackObject {
	if loadInterval < 0 || int64(loadInterval) >= time.Now().UnixNano() {
		var message = "Invalid loadInterval= " + strconv.Itoa(int(loadInterval))
		panic(message)
	}

	if loader == nil {
		var message = "loader= nil"
		panic(message)
	}

	var item = &BackObject{
		loadInterval: loadInterval,
		loader:       loader,
	}

	return item
}

func (item *BackObject) Get() interface{} {
	// init if this is the first time
	if atomic.LoadInt32(&item.done) == 0 {
		item.lock.Lock()
		defer item.lock.Unlock()

		if item.done == 0 {
			defer atomic.StoreInt32(&item.done, 1)
			if data, err := item.loader(); err == nil {
				item.data.Store(data)
			}

			// 每隔一段时间重新加载一次数据
			go func() {
				defer DumpIfPanic()
				var raw = int64(item.loadInterval)
				var step = raw >> 4
				for {
					// 将加载的间隔时间随机化，避免缓存雪崩
					var interval = raw + (rand.Int63n(step) << 1) - step
					time.Sleep(time.Duration(interval))
					if data, err := item.loader(); err == nil {
						item.data.Store(data)
					}
				}
			}()
		}
	}

	// 加载数据
	var data = item.data.Load()
	return data
}
