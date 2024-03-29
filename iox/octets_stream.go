package iox

import (
	"io"
)

/********************************************************************
created:    2023-06-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type OctetsStream struct {
	buffer   []byte // 有效内容为 buffer[position:len(buffer)]; read at &buffer[position], write at &buffer[len(buffer)]
	position int    // position只与read有关, 指向接下来要读取的位置; position<=len(buffer); 结构与iox.Buffer完全相同
}

func (my *OctetsStream) ReadBool() (bool, error) {
	var b, err = my.ReadByte()
	return b == 1, err
}

func (my *OctetsStream) ReadByte() (byte, error) {
	if my.position >= len(my.buffer) {
		return 0, ErrNotEnoughData
	}

	var result = my.buffer[my.position]
	my.position++
	return result, nil
}

func (my *OctetsStream) ReadInt16() (int16, error) {
	const readSize = 2
	if my.position+readSize > len(my.buffer) {
		return 0, ErrNotEnoughData
	}

	var b = my.buffer[my.position:]
	var result = int16(b[0]) | int16(b[1])<<8
	my.position += readSize
	return result, nil
}

func (my *OctetsStream) ReadInt32() (int32, error) {
	const readSize = 4
	if my.position+readSize > len(my.buffer) {
		return 0, ErrNotEnoughData
	}

	var b = my.buffer[my.position:]
	var result = int32(b[0]) | int32(b[1])<<8 | int32(b[2])<<16 | int32(b[3])<<24
	my.position += readSize
	return result, nil
}

func (my *OctetsStream) ReadInt64() (int64, error) {
	const readSize = 8
	if my.position+readSize > len(my.buffer) {
		return 0, ErrNotEnoughData
	}

	var b = my.buffer[my.position:]
	var result = int64(b[0]) | int64(b[1])<<8 | int64(b[2])<<16 | int64(b[3])<<24 | int64(b[4])<<32 | int64(b[5])<<40 | int64(b[6])<<48 | int64(b[7])<<56
	my.position += readSize
	return result, nil
}

func (my *OctetsStream) Read(buffer []byte) (int, error) {
	var readSize = len(buffer)
	if readSize == 0 {
		return 0, ErrInvalidArgument
	}

	var remainSize = len(my.buffer) - my.position
	if remainSize == 0 {
		return 0, nil
	}

	if readSize > remainSize {
		readSize = remainSize
	}

	copy(buffer, my.buffer[my.position:my.position+readSize])
	my.position += readSize
	return readSize, nil
}

func (my *OctetsStream) WriteBool(b bool) error {
	var d byte
	if b {
		d = 1
	}

	my.buffer = append(my.buffer, d)
	return nil
}

func (my *OctetsStream) WriteByte(b byte) error {
	my.buffer = append(my.buffer, b)
	return nil
}

func (my *OctetsStream) WriteInt16(d int16) error {
	my.buffer = append(my.buffer, byte(d), byte(d>>8))
	return nil
}

func (my *OctetsStream) WriteInt32(d int32) error {
	my.buffer = append(my.buffer, byte(d), byte(d>>8), byte(d>>16), byte(d>>24))
	return nil
}

func (my *OctetsStream) WriteInt64(d int64) error {
	my.buffer = append(my.buffer, byte(d), byte(d>>8), byte(d>>16), byte(d>>24), byte(d>>32), byte(d>>40), byte(d>>48), byte(d>>56))
	return nil
}

func (my *OctetsStream) Write(buffer []byte) error {
	var size = len(buffer)
	if size > 0 {
		my.buffer = append(my.buffer, buffer...)
	}

	return nil
}

func (my *OctetsStream) Len() int {
	return len(my.buffer)
}

func (my *OctetsStream) Position() int {
	return my.position
}

func (my *OctetsStream) Bytes() []byte {
	return my.buffer[my.position:]
}

func (my *OctetsStream) Tidy() {
	if my.position > 0 {
		copy(my.buffer, my.buffer[my.position:])
		my.buffer = my.buffer[:len(my.buffer)-my.position]
		my.position = 0
	}
}

func (my *OctetsStream) Reset() {
	my.buffer = my.buffer[:0]
	my.position = 0
}

func (my *OctetsStream) Seek(offset int64, whence int) (int64, error) {
	var num int64
	switch whence {
	case io.SeekStart:
		if offset < 0 {
			return 0, ErrInvalidArgument
		}
		num = 0
	case io.SeekCurrent:
		num = int64(my.position)
	case io.SeekEnd:
		num = int64(len(my.buffer))
	default:
		return 0, ErrInvalidArgument
	}

	num += offset
	if num < 0 {
		return 0, ErrInvalidArgument
	}

	my.position = int(num)
	return num, nil
}
