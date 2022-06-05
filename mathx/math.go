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

func ClampI32(i, a, b int32) int32 {
	if i < a {
		return a
	} else if a > b {
		return b
	} else {
		return i
	}
}

func ClampI64(i, a, b int64) int64 {
	if i < a {
		return a
	} else if a > b {
		return b
	} else {
		return i
	}
}

func ClampI(i, a, b int) int {
	if i < a {
		return a
	} else if a > b {
		return b
	} else {
		return i
	}
}

func ClampF32(f, a, b float32) float32 {
	if f < a {
		return a
	} else if a > b {
		return b
	} else {
		return f
	}
}

func ClampF64(f, a, b float64) float64 {
	if f < a {
		return a
	} else if a > b {
		return b
	} else {
		return f
	}
}
