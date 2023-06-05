package iox

import (
	"fmt"
	"io"
	"testing"
)

/********************************************************************
created:    2023-06-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

func TestOctetsStream_ReadByte(t *testing.T) {
	var stream = &OctetsStream{}
	for i := 0; i < 10; i++ {
		stream.WriteByte(byte(i))
	}

	stream.Seek(0, io.SeekStart)

	for {
		var b, err = stream.ReadByte()
		fmt.Printf("b=%v, err=%v\n", b, err)
		if err != nil {
			break
		}
	}
}

func TestOctetsStream_Read(t *testing.T) {
	var stream = &OctetsStream{}
	var count = 10
	for i := 0; i < count; i++ {
		stream.WriteByte(byte(i))
	}

	stream.Seek(0, io.SeekStart)

	var buffer = make([]byte, count)
	stream.Read(buffer, 0, count/2)
	stream.Read(buffer, 0, count/2)
	stream.Tidy()

	fmt.Println(buffer)
}
