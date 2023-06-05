package iox

import "errors"

/********************************************************************
created:    2023-06-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var ErrEmptyBuffer = errors.New("buffer is empty")
var ErrInvalidArgument = errors.New("invalid argument")
var ErrBad7BitInt = errors.New("bad 7bit int32")
var ErrNegativeSize = errors.New("size should not be negative")
var ErrNotEnoughData = errors.New("not enough data")
