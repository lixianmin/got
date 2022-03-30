package aesx

import "crypto/aes"

/********************************************************************
created:    2022-03-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type ICipher interface {
	Encrypt(input []byte) []byte
	Decrypt(input []byte) []byte
}

// NewCipher create a new cipher object.
func NewCipher(key []byte, options ...Option) ICipher {
	var block, err = aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	var args = arguments{
		cipherType:    cipherTypeCBC,
		initialVector: commonIV,
	}

	for _, opt := range options {
		opt(&args)
	}

	var cipher ICipher
	switch args.cipherType {
	case cipherTypeCFB:
		cipher = &cfbCipher{block: block, initialVector: args.initialVector}
	default:
		cipher = &cbcCipher{block: block, initialVector: args.initialVector}
	}

	return cipher
}
