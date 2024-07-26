package ants

import (
	"errors"
)

/********************************************************************
created:    2024-07-26
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

var errDiscard = errors.New("the task is discarded")

func IsDiscardError(err error) bool {
	return errors.Is(err, errDiscard)
}
