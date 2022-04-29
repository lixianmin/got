package taskx

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskEmpty struct{}

func (task taskEmpty) Do(args interface{}) error {
	return nil
}

func (task taskEmpty) Get1() interface{} {
	return nil
}

func (task taskEmpty) Get2() (interface{}, error) {
	return nil, nil
}
