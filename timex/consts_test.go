package timex

import (
	"fmt"
	"testing"
	"time"
)

/********************************************************************
created:    2021-03-04
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestFormatTime(t *testing.T) {
	var utc = time.Now().In(time.UTC)
	var local = time.Now()

	fmt.Printf("utc=%q, local=%q, utc-format=%q, local-format=%q\n", utc.Format(Layout), local.Format(Layout), FormatTime(utc), FormatTime(local))
}

func TestFormatDuration(t *testing.T) {
	var now = time.Now()
	var d = now.Sub(time.Date(1998, 10, 29, 12, 34, 56, 78, time.Local))
	fmt.Printf("d=%s\n", FormatDuration(d))
}
