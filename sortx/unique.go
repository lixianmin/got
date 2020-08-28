package sortx

/********************************************************************
created:    2020-08-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func UniqueInt(a []int) []int {
	var size = len(a)
	if size < 2 {
		return a
	}

	var j = 0
	for i := 1; i < size; i++ {
		if a[i] != a[j] {
			if j+1 != i {
				a[j+1] = a[i]
			}
			j++
		}
	}

	a = a[:j+1]
	return a
}

func UniqueString(a []string) []string {
	var size = len(a)
	if size < 2 {
		return a
	}

	var j = 0
	for i := 1; i < size; i++ {
		if a[i] != a[j] {
			if j+1 != i {
				a[j+1] = a[i]
			}
			j++
		}
	}

	a = a[:j+1]
	return a
}
