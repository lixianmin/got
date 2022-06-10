package ants

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type Task interface {
	Get1() interface{}
	Get2() (interface{}, error)
}
