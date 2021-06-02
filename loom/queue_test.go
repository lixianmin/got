package loom

import (
	"math/rand"
	"runtime"
	"strconv"
	"sync/atomic"
	"testing"
)

/********************************************************************
created:    2021-06-02
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestQueue(t *testing.T) {
	var q = NewQueue()

	count := 100
	for i := 0; i < count; i++ {
		q.Push(i)
	}

	for i := 0; i < count; i++ {
		v := q.Pop()
		if v == nil {
			t.Fatalf("got a nil value")
		}
		if v.(int) != i {
			t.Fatalf("expect %d but got %v", i, v)
		}
	}
}

func BenchmarkQueue(b *testing.B) {
	queues := map[string]*Queue{
		"lock-free queue": NewQueue(),
	}

	length := 1 << 12
	inputs := make([]int, length)
	for i := 0; i < length; i++ {
		inputs = append(inputs, rand.Int())
	}

	for _, cpus := range []int{4, 32, 1024} {
		runtime.GOMAXPROCS(cpus)
		for name, q := range queues {
			b.Run(name+"#"+strconv.Itoa(cpus), func(b *testing.B) {
				b.ResetTimer()

				var c int64
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						i := int(atomic.AddInt64(&c, 1)-1) % length
						v := inputs[i]
						if v >= 0 {
							q.Push(v)
						} else {
							q.Pop()
						}
					}
				})
			})
		}
	}
}
