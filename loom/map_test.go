package loom

//
//import (
//	"sync"
//	"testing"
//)
//
///********************************************************************
//created:    2020-07-13
//author:     lixianmin
//
//Copyright (C) - All Rights Reserved
//*********************************************************************/
//
//func TestMap_PutDifferentIntKeys(t *testing.T) {
//	var m Map
//
//	const key = 1029
//	m.Put(key, key)
//
//	if m.Get1(key) != key {
//		t.Fatalf("value should be %d", key)
//	}
//
//	if m.Get1(int16(key)) != key {
//		t.Fatalf("value should be %d", key)
//	}
//
//	if m.Get1(uint16(key)) != key {
//		t.Fatalf("value should be %d", key)
//	}
//
//	if m.Get1(int32(key)) != key {
//		t.Fatalf("value should be %d", key)
//	}
//
//	if m.Get1(uint32(key)) != key {
//		t.Fatalf("value should be %d", key)
//	}
//
//	if m.Get1(int64(key)) != key {
//		t.Fatalf("value should be %d", key)
//	}
//
//	if m.Get1(uint64(key)) != key {
//		t.Fatalf("value should be %d", key)
//	}
//}
//
//func TestMap_PutStringAndIntKeys(t *testing.T) {
//	var m Map
//	m.Put("hello", "world")
//
//	const intKey = 1029
//	m.ComputeIfAbsent(intKey, func(key any) any {
//		return key
//	})
//
//	if m.Size() != 2 {
//		t.Fatalf("size should be 2")
//	}
//
//	if m.Get1("hello") != "world" {
//		t.Fatalf("value should be world")
//	}
//
//	if m.Get1(intKey) != intKey {
//		t.Fatalf("value should be %d", intKey)
//	}
//}
//
//func TestMap_ComputeIfAbsent(t *testing.T) {
//	t.Parallel()
//	var m Map
//
//	const max = 1000
//	for i := 0; i < max; i++ {
//		m.Put(i, i)
//	}
//
//	const max2 = 2000
//	for i := max / 2; i < max2; i++ {
//		m.ComputeIfAbsent(i, func(key any) any {
//			return key.(int) * 2
//		})
//	}
//
//	if m.Get1(500) != 500 {
//		t.Fail()
//	}
//
//	if m.Get1(1500) != 3000 {
//		t.Fail()
//	}
//
//	if m.Size() != max2 {
//		t.Fail()
//	}
//}
//
//func TestMap_PutIfAbsent(t *testing.T) {
//	t.Parallel()
//	var m Map
//	const max = 1000
//
//	for i := 0; i < max; i++ {
//		m.Put(i, i)
//	}
//
//	for i := 0; i < max; i++ {
//		m.PutIfAbsent(i, i*2)
//	}
//
//	if m.Get1(5) != 5 {
//		t.Fail()
//	}
//
//	var input = 1
//	var ret = m.PutIfAbsent(max+1, input)
//	if ret != nil {
//		t.Fail()
//	}
//
//	ret = m.PutIfAbsent(max+1, input)
//	if ret != input {
//		t.Fail()
//	}
//}
//
//func TestMap_Remove(t *testing.T) {
//	t.Parallel()
//	var m Map
//	const max = 1000
//
//	for i := 0; i < max; i++ {
//		m.Put(i, i)
//	}
//
//	for i := 0; i < max; i++ {
//		m.Remove(i)
//	}
//
//	if m.Size() != 0 {
//		t.Fail()
//	}
//}
//
//func TestMap_Range(t *testing.T) {
//	t.Parallel()
//	var m Map
//	const max = 1000
//
//	m.Range(func(key any, value any) {
//
//	})
//
//	for i := 0; i < max; i++ {
//		m.Put(i, i)
//	}
//
//	var counter = 0
//	m.Range(func(key any, value any) {
//		counter += 1
//		//fmt.Println(key, value)
//	})
//
//	if counter != m.Size() {
//		t.Fail()
//	}
//}
//
//func TestMap_RangeDataRace(t *testing.T) {
//	var m Map
//
//	var wg sync.WaitGroup
//	wg.Add(2)
//
//	go func() {
//		for i := 0; i < 1000; i++ {
//			m.Put(i, i)
//		}
//
//		wg.Done()
//	}()
//
//	go func() {
//		for i := 0; i < 10; i++ {
//			m.Range(func(key interface{}, value interface{}) {
//
//			})
//		}
//		wg.Done()
//	}()
//
//	wg.Wait()
//}
//
////func TestMap_WithLock(t *testing.T) {
////	var m Map
////	m.Put(1, 1)
////
////	m.WithLock(1, func(table ShardTable) {
////		var v, ok = table[1]
////		fmt.Print(v, ok)
////	})
////}
