package convert

import (
	"fmt"
	"strconv"
	"unsafe"
)

/********************************************************************
created:    2018-09-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func ToString(v interface{}) string {
	switch v := v.(type) {
	case []byte:
		return *(*string)(unsafe.Pointer(&v))
	case int:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case string:
		return v
	default:
		return fmt.Sprintf("convert.ToString(): unexpected type for ToString, got type %T", v)
	}
}

func ToHuman(num uint64) string {
	if num >= 1073741824 {
		var v = float64(num) / 1073741824
		return fmt.Sprintf("%.2fG", v)
	} else if num >= 1048576 {
		var v = float64(num) / 1048576
		return fmt.Sprintf("%.2fM", v)
	} else if num >= 1024 {
		var v = float64(num) / 1024
		return fmt.Sprintf("%.2fK", v)
	} else {
		return fmt.Sprintf("%dB", num)
	}
}