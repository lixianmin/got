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
// 标准库中有一个sort.Reverse()方法，但是需要传入的sort.Interface类型; 这里实现的这个sortx.Reverse()传入slice即可
func Reverse(list interface{}) {
	var listValue = reflect.ValueOf(list)
	var count = listValue.Len()
	if count > 1 {
		var swapper = reflect.Swapper(listValue.Interface())
		for i, j := 0, count-1; i < j; i, j = i+1, j-1 {
			swapper(i, j)
		}
	}
}
