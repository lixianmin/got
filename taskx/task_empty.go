package taskx

/********************************************************************
created:    2020-08-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskEmpty struct{}

func (task taskEmpty) Do(args any) error {
	return nil
}

func (task taskEmpty) Get1() any {
	return nil
}

func (task taskEmpty) Get2() (any, error) {
	return nil, nil
}
