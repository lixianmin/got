package iox

import (
	"fmt"
	"github.com/lixianmin/got/convert"
	"testing"
)

/********************************************************************
created:    2022-01-11
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

// BytesToInt32 decode packet data length byte to int(Big end)
func BytesToInt(b []byte) int {
	var result = 0
	for _, v := range b {
		result = result<<8 + int(v)
	}

	return result
}

// IntToBytes encode packet data length to bytes(Big end)
func IntToBytes(n int) []byte {
	buf := make([]byte, 3)
	buf[0] = byte((n >> 16) & 0xFF)
	buf[1] = byte((n >> 8) & 0xFF)
	buf[2] = byte(n & 0xFF)
	return buf
}

func writeString(input *Buffer, s string) {
	var size = len(s)
	_, _ = input.Write(IntToBytes(size))
	_, _ = input.Write(convert.Bytes(s))
}

func readString(input *Buffer) string {
	var head = BytesToInt(input.Next(3))
	var slice = input.Next(head)
	var text = convert.String(slice)
	return text
}

func TestBuffer_Tidy(t *testing.T) {
	var input = &Buffer{}
	var count = 10
	for i := 0; i < count; i++ {
		writeString(input, fmt.Sprintf("hello:%d", i))
	}

	var i = 0
	for input.Len() > 0 {
		input.Tidy()
		var text = readString(input)
		println(text)

		if (i & 1) == 0 {
			writeString(input, "world")
			input.Tidy()
		}

		i++
	}
}
