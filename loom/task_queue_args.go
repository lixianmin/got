package loom

/********************************************************************
created:    2020-09-28
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type TaskQueueArgs struct {
	Size      int
	CloseChan chan struct{}
}

func (args *TaskQueueArgs) checkInit() {
	if args.Size <= 0 {
		args.Size = 8
	}

	if args.CloseChan == nil {
		args.CloseChan = make(chan struct{})
	}
}
