package sortx

import (
	"fmt"
	"testing"
)

/********************************************************************
created:    2020-07-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestReverseSlice(t *testing.T) {
	var list = []int{1, 2, 3, 4, 5, 6}
	var clonedList = make([]int, len(list))
	copy(clonedList, list)

	Reverse(list)
	for i := 0; i < len(list); i++ {
		if clonedList[i] != list[len(list)-i-1] {
			t.Fail()
		}
	}

	Reverse(list)
	fmt.Println(list)
}

func TestReverseNil(t *testing.T) {
	var list []int = nil
	Reverse(list)
}
