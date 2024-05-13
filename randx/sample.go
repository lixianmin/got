package randx

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
)

/********************************************************************
created:    2020-08-24
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type sampleHeapItem struct {
	ki    float64
	index int
}

type sampleHeap []sampleHeapItem

func (h *sampleHeap) Less(i, j int) bool {
	return (*h)[i].ki < (*h)[j].ki
}

func (h *sampleHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *sampleHeap) Len() int {
	return len(*h)
}

func (h *sampleHeap) Pop() (v any) {
	*h, v = (*h)[:h.Len()-1], (*h)[h.Len()-1]
	return
}

func (h *sampleHeap) Push(v any) {
	*h = append(*h, v.(sampleHeapItem))
}

func (h *sampleHeap) Get(index int) sampleHeapItem {
	return (*h)[index]
}

// 参考文献：http://lotabout.me/2018/Weighted-Random-Sampling
// 加权采样，返回索引下标
func WeightedSampling(sampleNum int, totalNum int, getWeight func(int) float64) []int {
	if totalNum < sampleNum || totalNum <= 0 {
		var message = fmt.Sprintf("invalid inputs: sampleNum=%d, totalNum=%d", sampleNum, totalNum)
		panic(message)
	}

	type Item struct {
		ki    float64
		index int
	}

	h := make(sampleHeap, sampleNum)
	for i := 0; i < totalNum; i++ {
		ui := rand.Float64()
		ki := math.Pow(ui, 1/getWeight(i))

		if h.Len() < sampleNum {
			heap.Push(&h, sampleHeapItem{ki: ki, index: i})
		} else if ki > h.Get(0).ki {
			heap.Push(&h, sampleHeapItem{ki: ki, index: i})
			if h.Len() > sampleNum {
				heap.Pop(&h)
			}
		}
	}

	var results = make([]int, sampleNum)
	for i := 0; i < sampleNum; i++ {
		results[i] = h.Get(i).index
	}

	return results
}
