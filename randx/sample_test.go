package randx

import (
	"fmt"
	"testing"
)

/********************************************************************
created:    2020-08-24
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 加权采样，返回索引下标
func TestWeightedSampling(t *testing.T) {
	var weights = []float64{3.0, 9.0, 2.0, 1.0, 1.0, 1.0, 2.0}
	var counter = make([]int, len(weights))

	for k := 0; k < 100000; k++ {
		var results = WeightedSampling(3, len(weights), func(i int) float64 {
			return weights[i]
		})

		for i := 0; i < len(results); i++ {
			counter[results[i]] ++
		}
	}

	fmt.Printf("%+v\n", counter)
}
