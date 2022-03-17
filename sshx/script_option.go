package sshx

/********************************************************************
created:    2021-11-25
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

const defaultName = "sshx.script"

type scriptOptions struct {
	name string
	//sha256     bool
}

type ScriptOption func(*scriptOptions)

// WithName script的文件名前缀, 默认值 sshx.script
func WithName(name string) ScriptOption {
	return func(options *scriptOptions) {
		if name != "" {
			options.name = name
		}
	}
}

// 加入这个参数的愿意, 是周期性生成含有时间参数的脚本, 可以复用同一个文件名. 但是, 发现如果直接使用文件名, 当script内容变化时, 不会被新内容覆盖, 这就起不到计划中的作用了
//// WithSha256 脚本文件名是否带sha256, 默认值 true
//func WithSha256(enable bool) ScriptOption {
//	return func(options *scriptOptions) {
//		options.sha256 = enable
//	}
//}
