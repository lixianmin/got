package timex

import (
	"time"
)

/********************************************************************
created:    2020-07-23
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const (
	TimeLayout = "2006-01-02 15:04:05"
	Day        = 24 * time.Hour
	Week       = 7 * Day
)

// 按照国人习惯的方式格式化了一下时间
func FormatTime(t time.Time) string {
	return t.Format(TimeLayout)
}