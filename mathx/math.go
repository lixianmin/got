package mathx

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func Min(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func Max(a int, b int) int {
	if a < b {
		return b
	}

	return a
}

func Clampf64(f float64, a float64, b float64) float64 {
	if f < a {
		return a
	} else if a > b {
		return b
	} else {
		return f
	}
}
