package taskx

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Task interface {
	Do(args any) error
	Get1() any
	Get2() (any, error)
}
