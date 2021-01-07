package sortx

import (
	"testing"
)

/********************************************************************
created:    2021-01-07
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestSliceBy(t *testing.T) {
	var keys = []int{3, 9, 7, 2, 100, 0, 4, 6}
	var size = len(keys)
	var values = make([]int, size)
	copy(values, keys)

	SliceBy(keys, values, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for i:= 0; i< len(keys); i++ {
		if keys[i] != values[i] {
			t.Fatalf("i=%d, keys[i]=%d, values[i]=%d", i, keys[i], values[i])
		}
	}
}