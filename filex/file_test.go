package filex

import (
	"fmt"
	"testing"
)

/********************************************************************
created:    2020-09-22
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestReadLines(t *testing.T) {
	_ = ReadLines("file_test.go", func(line string) bool {
		fmt.Println(line)
		return true
	})
}

func TestReadAllLines(t *testing.T) {
	var lines, _ = ReadAllLines("file_test.go")
	for _, line := range lines {
		fmt.Println(line)
	}
}