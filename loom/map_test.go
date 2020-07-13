package loom

import "testing"

/********************************************************************
created:    2020-07-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestMap_ComputeIfAbsent(t *testing.T) {
	var m Map

	for i := 0; i < 10; i++ {
		m.Put(i, i)
	}

	for i := 5; i < 15; i++ {
		m.ComputeIfAbsent(i, func(key interface{}) interface{} {
			return key.(int) * 2
		})
	}

	if m.Get(5) != 5 {
		t.Fail()
	}

	if m.Get(15) != 30 {
		t.Fail()
	}
}

func TestMap_Get(t *testing.T) {

}

func TestMap_Put(t *testing.T) {

}

func TestMap_PutIfAbsent(t *testing.T) {

}

func TestMap_Remove(t *testing.T) {

}

func TestMap_Size(t *testing.T) {

}
