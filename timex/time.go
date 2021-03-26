package timex

import "time"

/********************************************************************
created:    2020-07-23
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// 按照国人习惯的方式格式化了一下时间
func FormatTime(t time.Time) string {
	if t.Location() != time.Local {
		t = t.In(time.Local)
	}

	return t.Format(Layout)
}
