package sshx

/********************************************************************
created:    2021-11-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const defaultPrefix = "ssh."

type sshOptions struct {
	prefix string // script的前缀
	debug  bool   // 是否是debug模式
}

type SSHOption func(*sshOptions)

func WithPrefix(prefix string) SSHOption {
	return func(options *sshOptions) {
		if prefix != "" {
			options.prefix = prefix
		}
	}
}

func WithDebug() SSHOption {
	return func(options *sshOptions) {
		options.debug = true
	}
}