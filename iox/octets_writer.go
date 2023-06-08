package iox

import "github.com/lixianmin/got/convert"

/********************************************************************
created:    2023-06-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type OctetsWriter struct {
	stream *OctetsStream
}

func NewOctetsWriter(stream *OctetsStream) *OctetsWriter {
	var my = &OctetsWriter{stream: stream}
	return my
}

func (my *OctetsWriter) WriteBool(b bool) error {
	return my.stream.WriteBool(b)
}

func (my *OctetsWriter) WriteByte(b byte) error {
	return my.stream.WriteByte(b)
}

func (my *OctetsWriter) WriteInt16(d int16) error {
	return my.stream.WriteInt16(d)
}

func (my *OctetsWriter) WriteInt32(d int32) error {
	return my.stream.WriteInt32(d)
}

func (my *OctetsWriter) WriteInt64(d int64) error {
	return my.stream.WriteInt64(d)
}

func (my *OctetsWriter) WriteString(s string) error {
	var data = convert.Bytes(s)
	return my.WriteBytes(data)
}

func (my *OctetsWriter) WriteBytes(data []byte) error {
	var size = len(data)
	if err := my.Write7BitEncodedInt(int32(size)); err != nil {
		return err
	}

	return my.stream.Write(data)
}

// Write7BitEncodedInt 开启整数压缩
func (my *OctetsWriter) Write7BitEncodedInt(d int32) error {
	var num = uint32(d)
	for num > 127 {
		_ = my.stream.WriteByte(byte(num | 0xFFFFFF80))
		num >>= 7
	}

	return my.stream.WriteByte(byte(num))
}

func (my *OctetsWriter) Stream() *OctetsStream {
	return my.stream
}
