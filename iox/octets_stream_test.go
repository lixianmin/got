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

	var _, err = stream.ReadByte()
	if err != ErrNotEnoughData {
		t.Fail()
	}
}

func TestOctetsStream_Read(t *testing.T) {
	var stream = &OctetsStream{}
	const count = 10
	for i := 0; i < count; i++ {
		_ = stream.WriteByte(byte(i))
	}

	_, _ = stream.Seek(0, io.SeekStart)

	var buffer = make([]byte, count/2)
	var num, _ = stream.Read(buffer)
	if num != count/2 {
		t.Fail()
	}

	for i := 0; i < count/2; i++ {
		if buffer[i] != byte(i) {
			t.Fail()
		}
	}

	stream.Tidy()
	if stream.Len() != count/2 {
		t.Fail()
	}
}

func TestOctetsStream_Bytes(t *testing.T) {
	var stream = &OctetsStream{}
	stream.Bytes()
}

func TestOctetsStream_ReadInt32(t *testing.T) {
	var stream = &OctetsStream{}
	// stream里的数据为: 78, 97, 188, 0, 与charp中BinaryWriter中的是一样的
	_ = stream.WriteInt32(12345678)
	stream.Tidy()
}
