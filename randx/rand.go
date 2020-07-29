package randx

import (
	"math/rand"
	"time"
)

/********************************************************************
created:    2020-07-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// [from, to)
func Duration(from, to time.Duration) time.Duration {
	if from >= to {
		panic("from >= to")
	}

	return from + time.Duration(rand.Int63n(int64(to-from)))
}
