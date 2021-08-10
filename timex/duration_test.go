package timex

import (
	"fmt"
	"testing"
	"time"
)

/********************************************************************
created:    2021-08-10
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestFormatDuration(t *testing.T) {
	var now = time.Now()
	var d = now.Sub(time.Date(1998, 10, 29, 12, 34, 56, 78, time.Local))
	fmt.Printf("d=%s\n", FormatDuration(d))

	d = time.Hour*10 + time.Minute*9 + time.Second*8
	fmt.Printf("d=%s\n", FormatDuration(d))
}
