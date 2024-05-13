package std

/********************************************************************
created:    2021-04-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// Logger is used for logging formatted messages.
type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...any)
}

// LoggerFunc is a bridge between Logger and any third party logger
type LoggerFunc func(format string, args ...any)

func (f LoggerFunc) Printf(format string, args ...any) { f(format, args...) }
