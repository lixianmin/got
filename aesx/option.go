package aesx

/********************************************************************
created:    2022-03-29
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type arguments struct {
	cipherType    int
	initialVector []byte
}

type Option func(*arguments)

func WithCBC() Option {
	return func(args *arguments) {
		args.cipherType = cipherTypeCBC
	}
}

func WithCFB() Option {
	return func(args *arguments) {
		args.cipherType = cipherTypeCFB
	}
}

func WithInitialVector(initialVector []byte) Option {
	return func(args *arguments) {
		if len(initialVector) != 0 {
			args.initialVector = initialVector
		}
	}
}
