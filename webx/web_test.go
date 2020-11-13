package webx

import (
	"fmt"
	"testing"
	"time"
)

/********************************************************************
created:    2020-11-13
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestGet(t *testing.T) {
	var url = "http://www.baidu.com"
	var result, err = Get(url, WithTimeout(time.Second*2))
	fmt.Printf("result=%s, err=%q", result, err)
}
