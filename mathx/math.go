package mathx

/********************************************************************
created:    2019-06-19
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

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
