package loom

import (
	"fmt"
	"testing"
)

/********************************************************************
created:    2021-02-10
author:     lixianmin

这个为了later.go中stoppers排序做方法测试的
排序完成后：odd在左，even在右

Copyright (C) - All Rights Reserved
*********************************************************************/

func oddEvenSort(data []int) {
	if len(data) == 0 {
		return
	}

	var i, j = 0, len(data) - 1
	for {
		for i <= j && !isEven(data[i]) {
			i++
		}

		for i < j && isEven(data[j]) {
			j--
		}

		if i < j {
			data[i], data[j] = data[j], data[i]
		} else {
			fmt.Print(i, j, "\t")
			break
		}
	}

	fmt.Printf("%+v\t %+v\n", data, data[:i])
}

// 相当于IsStopped()方法
func isEven(d int) bool {
	return d&1 == 0
}

func TestOddEvenSort(t *testing.T) {
	oddEvenSort(nil)
	oddEvenSort([]int{1})
	oddEvenSort([]int{2})
	oddEvenSort([]int{1, 2, 3, 4, 5, 6, 7, 8})
	oddEvenSort([]int{1, 2, 3})
	oddEvenSort([]int{3, 2, 1})
	oddEvenSort([]int{1, 2})
	oddEvenSort([]int{2, 1})
	oddEvenSort([]int{2, 4, 6, 8})
	oddEvenSort([]int{1, 3, 5, 7})
}
