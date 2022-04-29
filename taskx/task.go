package taskx

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Task interface {
	Do(args interface{}) error
	Get1() interface{}
	Get2() (interface{}, error)
}
