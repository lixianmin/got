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

func TestOctetsReader_ReadByte(t *testing.T) {
	var stream = &OctetsStream{}
	var writer = NewOctetsWriter(stream)
	const count = 2
	for i := 0; i < count; i++ {
		writer.WriteByte(byte(i))
	}

	stream.Seek(0, io.SeekStart)
	var reader = NewOctetsReader(stream)
	for i := 0; i < count; i++ {
		fmt.Println(reader.ReadByte())
	}
}

func TestOctetsReader_ReadInt32(t *testing.T) {
	var stream = &OctetsStream{}
	var writer = NewOctetsWriter(stream)
	writer.WriteInt32(12345)
	writer.WriteInt32(678910)

	stream.Seek(0, io.SeekStart)
	var reader = NewOctetsReader(stream)
	fmt.Println(reader.ReadInt32())
	fmt.Println(reader.ReadInt32())
}

func TestOctetsReader_ReadString(t *testing.T) {
	var stream = &OctetsStream{}
	var writer = NewOctetsWriter(stream)
	writer.WriteString("中国")
	writer.WriteString("hello world")

	stream.Seek(0, io.SeekStart)
	var reader = NewOctetsReader(stream)
	fmt.Println(reader.ReadString())
	fmt.Println(reader.ReadString())
}

func TestOctetsReader_Read7BitEncodedInt(t *testing.T) {
	var stream = &OctetsStream{}
	var writer = NewOctetsWriter(stream)
	var reader = NewOctetsReader(stream)

	var d int32 = 126
	_ = writer.Write7BitEncodedInt(d)
	stream.Seek(0, io.SeekStart)
	var a, _ = reader.Read7BitEncodedInt()

	if a != d {
		t.Fail()
	}

	stream.Reset()
	d = 1234567
	_ = writer.Write7BitEncodedInt(d)
	stream.Seek(0, io.SeekStart)
	a, _ = reader.Read7BitEncodedInt()
	if a != d {
		t.Fail()
	}
}
