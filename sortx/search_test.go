package sortx

import "testing"

/********************************************************************
created:    2020-07-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestSearchHit(t *testing.T) {
	var list = []int {1, 3, 3, 3, 5, 7, 9, 9, 9, 11}

	var target = 9
	var index = Search(len(list), func(i int) bool {
		return list[i] < target
	}, func(i int) bool {
		return list[i] == target
	})

	// 找到的位置，应该是target第一次出现的索引下标
	if index != 6 {
		t.Fail()
	}
}

func TestSearchMiss(t *testing.T) {
	var list = []int {1, 3, 3, 5, 7, 9, 11}

	var target = 4
	var index = Search(len(list), func(i int) bool {
		return list[i] < target
	}, func(i int) bool {
		return list[i] == target
	})

	// 找不到的时候，返回的index应该是负数，且index的相反数应该是将target插入到有序列表中时它所在的目标下标
	index = ^index
	if index != 3 {
		t.Fail()
	}
}

func TestSearchHitReverse(t *testing.T) {
	var list = []int {11, 9, 7, 5, 4, 3, 1}

	var target = 4
	var index = Search(len(list), func(i int) bool {
		return list[i] > target
	}, func(i int) bool {
		return list[i] == target
	})

	if index != 4 {
		t.Fail()
	}
}
