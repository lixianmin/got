package loom

import "sync"

/********************************************************************
created:    2021-02-01
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ListPool struct {
	pool sync.Pool
}

func newPool(listSize int) *ListPool {
	var my = &ListPool{
		pool: sync.Pool{
			New: func() interface{} {
				return make([]interface{}, 0, listSize)
			},
		},
	}

	return my
}

func (my *ListPool) Get() []interface{} {
	return my.pool.Get().([]interface{})
}

func (my *ListPool) Put(list []interface{}) {
	const maxSize = 4096
	if list != nil && len(list) <= maxSize {
		list = list[:0]
		my.pool.Put(list)
	}
}
