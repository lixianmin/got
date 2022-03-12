package loom

/********************************************************************
created:    2020-03-12
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

import (
	"strconv"
	"sync"
	"testing"
)

const benchCount = 10000

func BenchmarkLoomMap_Put(b *testing.B) {
	b.StartTimer()
	var m Map

	const max = benchCount
	for i := 0; i < max; i++ {
		m.Put(i, i)
	}
}

func BenchmarkLoomMap_ComputeIfAbsent(b *testing.B) {
	b.StartTimer()
	var m Map

	const max = benchCount
	for i := 0; i < max; i++ {
		m.ComputeIfAbsent(i, func(key interface{}) interface{} {
			return key
		})
	}
}

func BenchmarkLoomMap_Get1(b *testing.B) {
	b.StopTimer()
	var m Map

	const max = benchCount
	for i := 0; i < max; i++ {
		m.Put(i, i)

		var s = strconv.Itoa(i)
		m.Put(s, s)
	}

	b.StartTimer()
	for i := 0; i < max*2; i++ {
		m.Get1(i)

		var s = strconv.Itoa(i)
		m.Get1(s)
	}
}

func BenchmarkSyncMap_Store(b *testing.B) {
	b.StartTimer()
	var m sync.Map

	const max = benchCount
	for i := 0; i < max; i++ {
		m.Store(i, i)
	}
}

func BenchmarkSyncMap_Load(b *testing.B) {
	b.StopTimer()
	var m sync.Map

	const max = benchCount
	for i := 0; i < max; i++ {
		m.Store(i, i)

		var s = strconv.Itoa(i)
		m.Store(s, s)
	}

	b.StartTimer()
	for i := 0; i < max*2; i++ {
		m.Load(i)

		var s = strconv.Itoa(i)
		m.Load(s)
	}
}
