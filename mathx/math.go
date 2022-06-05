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

func ClampI32(v, min, max int32) int32 {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}

func ClampI64(v, min, max int64) int64 {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}

func ClampI(v, min, max int) int {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}

func ClampF32(v, min, max float32) float32 {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}

func ClampF64(v, min, max float64) float64 {
	if v < min {
		return min
	} else if v > max {
		return max
	} else {
		return v
	}
}
