package iox

import (
	"encoding/binary"
	"io"
)

/********************************************************************
created:    2023-06-05
author:     lixianmin

Copyright (C) - All Rights Reserved
*********************************************************************/

type OctetsStream struct {
	buffer     []byte
	position   int
	length     int // 有效数据长度
	dirtyBytes int
}

func (my *OctetsStream) ReadByte() (byte, error) {
	if my.position >= my.length {
		return 0, ErrNotEnoughData
	}

	var result = my.buffer[my.position]
	my.position++
	return result, nil
}

func (my *OctetsStream) ReadInt32() (int32, error) {
	const size = 4
	if my.position+size > my.length {
		return 0, ErrNotEnoughData
	}

	var result = int32(binary.LittleEndian.Uint32(my.buffer[my.position:]))
	my.position += size
	return result, nil
}

func (my *OctetsStream) Read(buffer []byte, offset int, count int) (int, error) {
	var size = len(my.buffer)
	if size == 0 {
		return 0, ErrEmptyBuffer
	}

	if offset < 0 || count < 0 {
		return 0, ErrInvalidArgument
	}

	if size-offset < count {
		return 0, ErrInvalidArgument
	}

	if my.position >= my.length || count == 0 {
		return 0, nil
	}

	var remaining = int(my.length - my.position)
	if count > remaining {
		count = remaining
	}

	copy(buffer[offset:], my.buffer[my.position:my.position+count])
	my.position += count
	return count, nil
}

func (my *OctetsStream) WriteByte(b byte) error {
	if my.position >= my.length {
		my.expand(my.position + 1)
		my.length = my.position + 1
	}

	my.buffer[my.position] = b
	my.position++
	return nil
}

func (my *OctetsStream) WriteInt32(d int32) error {
	const size = 4
	var targetLength = my.position + size
	if targetLength > my.length {
		my.expand(targetLength)
	}

	binary.LittleEndian.PutUint32(my.buffer[my.position:], uint32(d))
	my.position += size
	if my.position >= my.length {
		my.length = my.position
	}
	return nil
}

func (my *OctetsStream) Write(buffer []byte, offset int, count int) error {
	if buffer == nil {
		return ErrEmptyBuffer
	}

	if offset < 0 || count < 0 {
		return ErrInvalidArgument
	}

	if len(my.buffer)-offset < count {
		return ErrInvalidArgument
	}

	var targetLength = my.position + count
	if targetLength > my.length {
		my.expand(targetLength)
	}

	copy(my.buffer[my.position:], buffer[offset:offset+count])
	my.position += count

	if my.position >= my.length {
		my.length = my.position
	}

	return nil
}

func (my *OctetsStream) expand(nextSize int) {
	var lastSize = cap(my.buffer)
	if nextSize > lastSize {
		const minSize = 16
		if nextSize < minSize {
			nextSize = minSize
		} else if nextSize < (lastSize << 1) {
			nextSize = lastSize << 1
		}

		var array = make([]byte, nextSize)
		if my.length > 0 {
			copy(array, my.buffer[:my.length])
		}

		my.dirtyBytes = 0
		my.buffer = array
	} else if my.dirtyBytes > 0 {
		var size = my.dirtyBytes
		for i := 0; i < size; i++ {
			my.buffer[i] = 0
		}
		my.dirtyBytes = 0
	}
}

func (my *OctetsStream) SetLength(length int) error {
	if length < 0 || length > cap(my.buffer) {
		return ErrInvalidArgument
	}

	var num = length
	if num > my.length {
		my.expand(num)
	} else if num < my.length {
		my.dirtyBytes += my.length - num
	}

	my.length = num
	if my.position > my.length {
		my.position = my.length
	}

	return nil
}

func (my *OctetsStream) GetLength() int {
	return my.length
}

func (my *OctetsStream) GetPosition() int {
	return my.position
}

func (my *OctetsStream) GetBytes() []byte {
	var remainSize = my.length - my.position
	return my.buffer[my.position:remainSize]
}

func (my *OctetsStream) Tidy() error {
	var count = my.length - my.position
	copy(my.buffer, my.buffer[my.position:my.position+count])

	my.position = 0
	return my.SetLength(count)
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
		num = int64(my.length)
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
