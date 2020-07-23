package mathx

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func MinInt32(a int32, b int32) int32 {
	if a > b {
		return b
	}

	return a
}

func MinInt64(a int64, b int64) int64 {
	if a > b {
		return b
	}

	return a
}

func MinInt(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func MaxInt(a int, b int) int {
	if a < b {
		return b
	}

	return a
}

func MaxInt32(a int32, b int32) int32 {
	if a < b {
		return b
	}
	
	return a
}

func MaxInt64(a int64, b int64) int64 {
	if a < b {
		return b
	}

	return a
}

func ClampInt32(i int32, a int32, b int32) int32 {
	if i < a {
		return a
	} else if a > b {
		return b
	} else {
		return i
	}
}

func ClampInt64(i int64, a int64, b int64) int64 {
	if i < a {
		return a
	} else if a > b {
		return b
	} else {
		return i
	}
}

func ClampInt(i int, a int, b int) int {
	if i < a {
		return a
	} else if a > b {
		return b
	} else {
		return i
	}
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
