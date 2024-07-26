package ants

import "context"

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Task interface {
	Get1() any
	Get2() (any, error)
	Err() error
	run(ctx context.Context)
}
