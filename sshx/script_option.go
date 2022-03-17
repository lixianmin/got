package sshx

/********************************************************************
created:    2021-11-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const defaultScriptName = "sshx.script"

type scriptOptions struct {
	scriptName string
	sha256     bool
}

type ScriptOption func(*scriptOptions)

// WithScriptName script的文件名前缀, 默认值 sshx.script
func WithScriptName(name string) ScriptOption {
	return func(options *scriptOptions) {
		if name != "" {
			options.scriptName = name
		}
	}
}

// WithSha256 脚本文件名是否带sha256, 默认值 true
func WithSha256(enable bool) ScriptOption {
	return func(options *scriptOptions) {
		options.sha256 = enable
	}
}
