package timex

import "time"

/********************************************************************
created:    2020-07-23
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// FormatTime 按照国人习惯的方式格式化了一下时间
func FormatTime(t time.Time) string {
	if t.Location() != time.Local {
		t = t.In(time.Local)
	}

	return t.Format(Layout)
}

// Midnight 把时间对齐到午夜
func Midnight(t time.Time) time.Time {
	var year, month, day = t.Date()
	var midnight = time.Date(year, month, day, 0, 0, 0, 0, t.Location())
	return midnight
}
