package aesx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

/********************************************************************
created:    2020-07-15
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type cbcCipher struct {
	block         cipher.Block
	initialVector []byte
}

func (my *cbcCipher) Encrypt(input []byte) []byte {
	var encrypt = cipher.NewCBCEncrypter(my.block, my.initialVector)
	input = pkcs5Padding(input, aes.BlockSize)
	var output = make([]byte, len(input))
	encrypt.CryptBlocks(output, input)
	return output
}

func (my *cbcCipher) Decrypt(input []byte) []byte {
	var decrypt = cipher.NewCBCDecrypter(my.block, my.initialVector)
	var output = make([]byte, len(input))
	decrypt.CryptBlocks(output, input)
	output = pkcs5Trimming(output)
	return output
}

// padding for fulfill block size requirement
func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	var padding = blockSize - len(ciphertext)%blockSize
	var padText = bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// trim padding tail
func pkcs5Trimming(encrypt []byte) []byte {
	var padding = encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)]
}
