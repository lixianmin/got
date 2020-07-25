package sortx

import (
	"reflect"
)

/********************************************************************
created:    2020-07-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func Reverse(list interface{}) {
	var listValue = reflect.Indirect(reflect.ValueOf(list))
	var kind = listValue.Kind()
	if kind != reflect.Slice {
		panic("not a slice")
	}

	var count = listValue.Len()
	if count <= 1 {
		return
	}

	var swapper = reflect.Swapper(listValue.Interface())
	for i, j := 0, count-1; i < j; i, j = i+1, j-1 {
		swapper(i, j)
	}
}
