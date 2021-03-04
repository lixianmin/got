package timex

import (
	"testing"
	"time"
)

/********************************************************************
created:    2021-03-04
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestLayout(t *testing.T) {
	var s = time.Now().Format(Layout)
	println(s)
}