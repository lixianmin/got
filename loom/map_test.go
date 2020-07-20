package loom

import (
	"testing"
)

/********************************************************************
created:    2020-07-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestMap_ComputeIfAbsent(t *testing.T) {
	t.Parallel()
	var m Map

	const max = 1000
	for i := 0; i < max; i++ {
		m.Put(i, i)
	}

	const max2 = 2000
	for i := max / 2; i < max2; i++ {
		m.ComputeIfAbsent(i, func(key interface{}) interface{} {
			return key.(int) * 2
		})
	}

	if m.Get1(500) != 500 {
		t.Fail()
	}

	if m.Get1(1500) != 3000 {
		t.Fail()
	}

	if m.Size() != max2 {
		t.Fail()
	}
}

func TestMap_PutIfAbsent(t *testing.T) {
	t.Parallel()
	var m Map
	const max = 1000

	for i := 0; i < max; i++ {
		m.Put(i, i)
	}

	for i := 0; i < max; i++ {
		m.PutIfAbsent(i, i*2)
	}

	if m.Get1(5) != 5 {
		t.Fail()
	}

	var input = 1
	var ret = m.PutIfAbsent(max+1, input)
	if ret != nil {
		t.Fail()
	}

	ret = m.PutIfAbsent(max+1, input)
	if ret != input {
		t.Fail()
	}
}

func TestMap_Remove(t *testing.T) {
	t.Parallel()
	var m Map
	const max = 1000

	for i := 0; i < max; i++ {
		m.Put(i, i)
	}

	for i := 0; i < max; i++ {
		m.Remove(i)
	}

	if m.Size() != 0 {
		t.Fail()
	}
}

func TestMap_Range(t *testing.T) {
	t.Parallel()
	var m Map
	const max = 1000

	m.Range(func(key interface{}, value interface{}) {
		
	})

	for i := 0; i < max; i++ {
		m.Put(i, i)
	}

	var counter = 0
	m.Range(func(key interface{}, value interface{}) {
		counter += 1
		//fmt.Println(key, value)
	})

	if counter != m.Size() {
		t.Fail()
	}
}
