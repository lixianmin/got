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
		_ = wc.Close(nil)
	}()

	<-wc.C()
}

func TestWaitClose_Close(t *testing.T) {
	t.Parallel()
	var wc WaitClose

	// 即使未初始化过的wc，也应该可以正常调用结束
	var f = func() {
		_ = wc.Close(func() error {
			println("closed once")
			return nil
		})
	}

	go f()
	go f()
	go f()

	time.Sleep(time.Second)
}

func TestWaitClose_Close2(t *testing.T) {
	t.Parallel()
	var wc WaitClose

	wc.WaitUtil(time.Millisecond)

	var f = func() {
		_ = wc.Close(func() error {
			panic("hello")
			return nil
		})
		_ = wc.Close(nil)
	}

	go f()
	go f()
	go f()
}

func TestWaitClose_WaitUtil_Direct(t *testing.T) {
	var wc WaitClose
	wc.Close(nil)
	wc.WaitUtil(time.Second)

	var wc2 WaitClose
	wc2.Close(nil)
	wc.WaitUtil(time.Second)
}

func TestWaitClose_WaitUtil_afterInited(t *testing.T) {
	var wc WaitClose
	wc.C()
	wc.WaitUtil(time.Second)
}

func TestWaitClose_WaitUtil_Closed(t *testing.T) {
	var wc WaitClose
	wc.C()
	wc.Close(nil)
	wc.WaitUtil(time.Second)
}
