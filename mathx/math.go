package mathx

import (
	"cmp"
)

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func Clamp[T cmp.Ordered](v, min, max T) T {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}
