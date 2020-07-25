package sortx

import (
	"reflect"
)

/********************************************************************
created:    2020-07-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 通用逆序算法
// 目标：将一个slice逆序
// 算法实现参考：https://stackoverflow.com/questions/54858529/golang-reverse-a-arbitrary-slice
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
