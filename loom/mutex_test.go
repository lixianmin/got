package loom

import (
	"fmt"
	"testing"
	"time"
)

/********************************************************************
created:    2020-10-21
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestMutex_TryLock(t *testing.T) {
	var m = Mutex{}
	var locked = m.TryLock()

	if locked {
		defer m.Unlock()
	}

	go func() {
		m.Lock()
		defer m.Unlock()
		fmt.Println("locked")
	}()

	fmt.Printf("try lock: locked=%v\n", locked)
	time.Sleep(time.Second)
}

func TestMutex_TryLock2(t *testing.T) {
	var m = Mutex{}
	go func() {
		m.Lock()
		defer m.Unlock()
		fmt.Println("locked")
		time.Sleep(2 * time.Second)
	}()

	time.Sleep(time.Second)
	var locked = m.TryLock()

	if locked {
		defer m.Unlock()
	}

	fmt.Printf("try lock: locked=%v\n", locked)
	time.Sleep(time.Second)
}
