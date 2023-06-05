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
	buffer       []byte
	position     int32
	length       int32
	capacity     int32
	dirtyBytes   int32
	initialIndex int32
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

	copy(buffer[offset:], my.buffer[my.position:my.position+int32(count)])
	my.position += int32(count)
	return count, nil
}

func (my *OctetsStream) WriteByte(b byte) error {
	if my.position >= my.length {
		if err := my.expand(my.position + 1); err != nil {
			return err
		}
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
		if err := my.expand(targetLength); err != nil {
			return err
		}
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

	var targetLength = my.position + int32(count)
	if targetLength > my.length {
		if err := my.expand(targetLength); err != nil {
			return err
		}
	}

	copy(my.buffer[my.position:], buffer[offset:offset+count])
	my.position += int32(count)

	if my.position >= my.length {
		my.length = my.position
	}

	return nil
}

func (my *OctetsStream) expand(newSize int32) error {
	if newSize > my.capacity {
		var num = newSize
		const minSize = 32
		if num < minSize {
			num = minSize
		} else if num < (my.capacity << 1) {
			num = my.capacity << 1
		}

		return my.SetCapacity(int(num))
	} else if my.dirtyBytes > 0 {
		var size = int(my.dirtyBytes)
		for i := 0; i < size; i++ {
			my.buffer[i] = 0
		}
		my.dirtyBytes = 0
	}

	return nil
}

func (my *OctetsStream) SetCapacity(capacity int) error {
	var c = int32(capacity)
	if c != my.capacity {
		if c < 0 || c < my.length {
			return ErrInvalidArgument
		}

		if capacity != len(my.buffer) {
			var array = make([]byte, capacity)
			if my.length > 0 {
				copy(array, my.buffer[:my.length])
			}

			my.dirtyBytes = 0
			my.buffer = array
			my.capacity = c
		}
	}

	return nil
}

func (my *OctetsStream) SetLength(length int) error {
	if length < 0 || int32(length) > my.capacity {
		return ErrInvalidArgument
	}

	var num = int32(length) + my.initialIndex
	if num > my.length {
		if err := my.expand(num); err != nil {
			return err
		}
	} else if num < my.length {
		my.dirtyBytes += my.length - num
	}

	my.length = num
	if my.position > my.length {
		my.position = my.length
	}

	return nil
}

func (my *OctetsStream) Tidy() error {
	var count = my.length - my.position
	copy(my.buffer[my.initialIndex:], my.buffer[my.position:my.position+count])

	my.position = my.initialIndex
	return my.SetLength(int(count))
}

func (my *OctetsStream) Seek(offset int64, whence int) (int64, error) {
	var num int32
	switch whence {
	case io.SeekStart:
		if offset < 0 {
			return 0, ErrInvalidArgument
		}
		num = my.initialIndex
	case io.SeekCurrent:
		num = my.position
	case io.SeekEnd:
		num = my.length
	default:
		return 0, ErrInvalidArgument
	}

	num += int32(offset)
	if num < my.initialIndex {
		return 0, ErrInvalidArgument
	}

	my.position = num
	return int64(num), nil
}
