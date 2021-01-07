package sortx

import (
	"github.com/lixianmin/got/mathx"
	"reflect"
)

/********************************************************************
created:    2021-01-07
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 将keys和values同步排序
func SliceBy(keys interface{}, values interface{}, less func(i, j int) bool) {
	var keyValues = reflect.ValueOf(keys)
	var valValues = reflect.ValueOf(values)

	var length = mathx.MinInt(keyValues.Len(), valValues.Len())
	if length <= 1 {
		return
	}

	var keySwapper = reflect.Swapper(keys)
	var valSwapper = reflect.Swapper(values)
	var swapper = func(i int, j int) {
		keySwapper(i, j)
		valSwapper(i, j)
	}

	quickSort_func(lessSwap{less, swapper}, 0, length, maxDepth(length))
}

// lessSwap is a pair of Less and Swap function for use with the
// auto-generated func-optimized variant of sort.go in
// zfuncversion.go.
type lessSwap struct {
	Less func(i, j int) bool
	Swap func(i, j int)
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}
