package std

/********************************************************************
created:    2020-10-14
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

import "reflect"

// IsNil 方法中传入指针时，直接使用 if i == nil {} 是无法判断是不为nil的
func IsNil(i any) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}
