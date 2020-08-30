package loom

import (
	"testing"
	"time"
)

func TestWaitClose_Chan(t *testing.T) {
	t.Parallel()
	var wc WaitClose

	go func() {
		for {
			select {
			case <-wc.C():
				println("hello")
				return
			}
		}
	}()

	var ticker = time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				println(time.Now().String())
			case <-wc.C():
				println("world")
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-wc.C():
				println("pet")
				return
			}
		}
	}()

	time.Sleep(time.Second * 2)
	wc.Close(nil)
	time.Sleep(time.Second * 100)
}
