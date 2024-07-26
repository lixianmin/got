package ants

import (
	"context"
)

/********************************************************************
created:    2024-07-26
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type taskDiscard struct {
}

func newTaskDiscard() *taskDiscard {
	var my = &taskDiscard{}
	return my
}

func (my *taskDiscard) Get1() any {
	return nil
}

func (my *taskDiscard) Get2() (any, error) {
	return nil, errDiscard
}

func (my *taskDiscard) Error() error {
	return errDiscard
}

func (my *taskDiscard) run(ctx context.Context) {
}
