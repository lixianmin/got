package iox

import (
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
	const count = 10
	for i := 0; i < count; i++ {
		_ = stream.WriteByte(byte(i))
	}

	_, _ = stream.Seek(0, io.SeekStart)

	for i := 0; i < count; i++ {
		var b, _ = stream.ReadByte()
		if int(b) != i {
			t.Fail()
		}
	}
}

func TestOctetsStream_Read(t *testing.T) {
	var stream = &OctetsStream{}
	const count = 10
	for i := 0; i < count; i++ {
		_ = stream.WriteByte(byte(i))
	}

	_, _ = stream.Seek(0, io.SeekStart)

	var buffer = make([]byte, count)
	var num, _ = stream.Read(buffer, 0, count/2)
	if num != count/2 {
		t.Fail()
	}

	for i := 0; i < count/2; i++ {
		if buffer[i] != byte(i) {
			t.Fail()
		}
	}

	_ = stream.Tidy()
	if stream.GetLength() != count/2 {
		t.Fail()
	}
}
