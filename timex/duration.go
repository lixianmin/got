package timex

import (
	"github.com/lixianmin/got/convert"
	"strconv"
	"time"
)

/********************************************************************
created:    2021-03-26
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func FormatDuration(d time.Duration) string {
	var buf = make([]byte, 0, 32)

	if d >= Day {
		var days = d / Day
		d -= days * Day
		buf = strconv.AppendUint(buf, uint64(days), 10)
		buf = append(buf, 'd')
	}

	if d >= time.Hour {
		var hours = d / time.Hour
		d -= hours * time.Hour
		buf = strconv.AppendUint(buf, uint64(hours), 10)
		buf = append(buf, 'h')
	}

	if d >= time.Minute {
		var minutes = d / time.Minute
		d -= minutes * time.Minute
		buf = strconv.AppendUint(buf, uint64(minutes), 10)
		buf = append(buf, 'm')
	}

	var seconds = float64(d) / float64(time.Second)
	buf = strconv.AppendFloat(buf, seconds, 'f', 3, 64)
	buf = append(buf, 's')

	return convert.String(buf)
}
