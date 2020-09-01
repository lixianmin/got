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

	go func() {
		time.Sleep(2 * time.Second)
		wc.Close(nil)
	}()

	<-wc.C()
}

func TestWaitClose_Close(t *testing.T) {
	t.Parallel()
	var wc WaitClose

	wc.WaitUtil(time.Millisecond)

	var f = func() {
		wc.Close(func() {
			println("closed once")
		})
	}

	go f()
	go f()
	go f()
}

func TestWaitClose_Close2(t *testing.T) {
	t.Parallel()
	var wc WaitClose

	wc.WaitUtil(time.Millisecond)

	var f = func() {
		wc.Close(func() {
			panic("hello")
		})
		wc.Close(nil)
	}

	go f()
	go f()
	go f()
}
