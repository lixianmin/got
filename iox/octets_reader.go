package iox

import "github.com/lixianmin/got/convert"

/********************************************************************
created:    2023-06-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type OctetsReader struct {
	stream *OctetsStream
}

func NewOctetsReader(stream *OctetsStream) *OctetsReader {
	var my = &OctetsReader{stream: stream}
	return my
}

func (my *OctetsReader) ReadByte() (byte, error) {
	return my.stream.ReadByte()
}

func (my *OctetsReader) ReadInt32() (int32, error) {
	return my.stream.ReadInt32()
}

func (my *OctetsReader) ReadString() (string, error) {
	var size, err = my.Read7BitEncodedInt()
	if err != nil {
		return "", err
	}

	if size < 0 {
		return "", ErrNegativeSize
	}

	if size == 0 {
		return "", nil
	}

	var data = make([]byte, size)
	var num, err2 = my.stream.Read(data)
	if err2 != nil {
		return "", err2
	}

	if int32(num) != size {
		return "", ErrNotEnoughData
	}

	var result = convert.String(data)
	return result, nil
}

func (my *OctetsReader) Read7BitEncodedInt() (int32, error) {
	var num uint32 = 0
	for i := 0; i < 28; i += 7 {
		var b, err = my.ReadByte()
		if err != nil {
			return 0, err
		}

		num |= (uint32)((b & 0x7F) << i)
		if b < 127 {
			return int32(num), nil
		}
	}

	var b, err = my.ReadByte()
	if err != nil {
		return 0, err
	}

	if b > 15 {
		return 0, ErrBad7BitInt
	}

	return int32(num) | (int32(b) << 28), nil
}

func (my *OctetsReader) GetStream() *OctetsStream {
	return my.stream
}
