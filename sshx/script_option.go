package sshx

/********************************************************************
created:    2021-11-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const defaultPrefix = "ssh."

type scriptOptions struct {
	prefix string // script的前缀
}

type ScriptOption func(*scriptOptions)

func WithPrefix(prefix string) ScriptOption {
	return func(options *scriptOptions) {
		if prefix != "" {
			options.prefix = prefix
		}
	}
}