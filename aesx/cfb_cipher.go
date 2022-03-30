package aesx

import (
	"crypto/cipher"
)

/********************************************************************
created:    2022-03-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type cfbCipher struct {
	block         cipher.Block
	initialVector []byte
}

func (my *cfbCipher) Encrypt(input []byte) []byte {
	var encrypt = cipher.NewCFBEncrypter(my.block, my.initialVector)
	var output = make([]byte, len(input))
	encrypt.XORKeyStream(output, input)
	return output
}

func (my *cfbCipher) Decrypt(input []byte) []byte {
	var decrypt = cipher.NewCFBDecrypter(my.block, my.initialVector)
	var output = make([]byte, len(input))
	decrypt.XORKeyStream(output, input)
	return output
}
