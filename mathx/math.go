package mathx

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func MinI32(a int32, b int32) int32 {
	if a > b {
		return b
	}

	return a
}

func MinI64(a int64, b int64) int64 {
	if a > b {
		return b
	}

	return a
}

func MinI(a int, b int) int {
	if a > b {
		return b
	}

	return a
}

func MaxI(a int, b int) int {
	if a < b {
		return b
	}

	return a
}

func MaxI32(a int32, b int32) int32 {
	if a < b {
		return b
	}

	return a
}

func MaxI64(a int64, b int64) int64 {
	if a < b {
		return b
	}

	return a
}

func ClampI32(i, min, max int32) int32 {
	if i < min {
		return min
	} else if min > max {
		return max
	} else {
		return i
	}
}

func ClampI64(i, min, max int64) int64 {
	if i < min {
		return min
	} else if min > max {
		return max
	} else {
		return i
	}
}

func ClampI(i, min, max int) int {
	if i < min {
		return min
	} else if min > max {
		return max
	} else {
		return i
	}
}

func ClampF32(f, min, max float32) float32 {
	if f < min {
		return min
	} else if min > max {
		return max
	} else {
		return f
	}
}

func ClampF64(f, min, max float64) float64 {
	if f < min {
		return min
	} else if min > max {
		return max
	} else {
		return f
	}
}
