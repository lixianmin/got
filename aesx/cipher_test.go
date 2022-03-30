package aesx

import (
	"fmt"
	"testing"
)

/********************************************************************
created:    2022-03-30
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestCbcCipher_Encrypt(t *testing.T) {
	var cipher = NewCipher([]byte("this-must-be-of-length-16-24-32."), WithInitialVector([]byte("equal-block-size")), WithCBC())
	var input = "hello world"
	var data = cipher.Encrypt([]byte(input))
	var output = string(cipher.Decrypt(data))

	if input != output {
		t.Fail()
	}

	fmt.Printf("input=%s, output=%s", input, output)
}

func TestCfbCipher_Encrypt(t *testing.T) {
	var cipher = NewCipher([]byte("this-must-be-of-length-16-24-32."), WithInitialVector([]byte("1234567890123456")), WithCBC())
	var input = "hello world"
	var data = cipher.Encrypt([]byte(input))
	var output = string(cipher.Decrypt(data))

	if input != output {
		t.Fail()
	}

	fmt.Printf("input=%s, output=%s", input, output)
}
