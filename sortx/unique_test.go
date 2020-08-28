package sortx

import "testing"

/********************************************************************
created:    2020-08-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestUniqueInt(t *testing.T) {
	var a = []int{1, 2, 2, 2, 4, 4, 6, 6, 6, 6, 6, 6}
	var b = []int{1, 2, 4, 6}
	a = UniqueInt(a)

	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			t.Fail()
		}
	}
}

func TestUniqueInt2(t *testing.T) {
	var a = []int{1, 2, 3, 4, 5, 6}
	var b = []int{1, 2, 3, 4, 5, 6}
	a = UniqueInt(a)

	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			t.Fail()
		}
	}
}
