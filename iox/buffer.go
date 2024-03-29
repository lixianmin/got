package iox

import (
	"errors"
	"io"
)

/********************************************************************
created:    2020-12-07
author:     lixianmin

参考 bytes.Buffer改遍而来，目前似乎只是加了一个Tidy()的方法

Copyright (C) - All Rights Reserved
*********************************************************************/

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("iox.Buffer: too large")

// var errNegativeRead = errors.New("iox.Buffer: reader returned negative count from Read")
var errInvalidSeek = errors.New("iox.Buffer: invalid seek")

const maxInt = int(^uint(0) >> 1)

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	buf []byte // contents are the bytes buf[off : len(buf)]
	off int    // read at &buf[off], write at &buf[len(buf)]

	// 通过buffer.Bytes()预取并计算，可以避免make checkpoint相关的逻辑
	//checkpointBuffer []byte
	//checkpointOffset int
}

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like Read, Write, Reset, or Truncate).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (b *Buffer) Bytes() []byte { return b.buf[b.off:] }

// String returns the contents of the unread portion of the buffer
// as a string. If the Buffer is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the strings.Builder type.
func (b *Buffer) String() string {
	if b == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(b.buf[b.off:])
}

// empty reports whether the unread portion of the buffer is empty.
func (b *Buffer) empty() bool { return len(b.buf) <= b.off }

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (b *Buffer) Len() int { return len(b.buf) - b.off }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (b *Buffer) Cap() int { return cap(b.buf) }

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as Truncate(0).
func (b *Buffer) Reset() {
	b.buf = b.buf[:0] // even if b.buf=nil, it can run without panic
	b.off = 0
}

func (b *Buffer) Seek(offset int64, whence int) (ret int64, err error) {
	if whence >= io.SeekStart && whence <= io.SeekEnd {
		var off = int(offset)
		var next = off

		if whence == io.SeekCurrent {
			next += b.off
		} else if whence == io.SeekEnd {
			next += len(b.buf)
		}

		if next >= 0 && next <= len(b.buf) {
			b.off = next
			return int64(next), nil
		}
	}

	return 0, errInvalidSeek
}

// tryGrowByReslice is an inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (b *Buffer) tryGrowByReslice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}

	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) grow(n int) int {
	m := b.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && b.off != 0 {
		b.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := b.tryGrowByReslice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(b.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(b.buf, b.buf[b.off:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := makeSlice(2*c + n)
		copy(buf, b.buf[b.off:])
		b.buf = buf
	}

	// Restore b.off and len(b.buf).
	b.off = 0
	b.buf = b.buf[:m+n]
	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
func (b *Buffer) Write(p []byte) (n int, err error) {
	m, ok := b.tryGrowByReslice(len(p))
	if !ok {
		m = b.grow(len(p))
	}
	return copy(b.buf[m:], p), nil
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
func (b *Buffer) Read(p []byte) (n int, err error) {
	if b.empty() {
		// 读的时候不调用Reset()方法，这样才有机会通过SetOffset()重置进度
		//// Buffer is empty, reset to recover space.
		//b.Reset()
		if len(p) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}

	n = copy(p, b.buf[b.off:])
	b.off += n
	return n, nil
}

// ReadOnce 主要用于读一次网络数据这种永不结束的流
func (b *Buffer) ReadOnce(reader io.Reader, buf []byte) (int, error) {
	// 类似于ReadFrom()，但ReadFrom()的结束条件为读到io.EOF或err
	// 为什么不采用循环读的方式？因为有可能对方一直发大规模的数据，永不结束，这样不但把iox.Buffer打死了，读过程还结束不了
	//
	// 本方法是从其它的reader读数据到本buffer中，只修改b.buf，不修改b.off
	var num, err = reader.Read(buf)
	if err != nil {
		return 0, err
	}

	_, _ = b.Write(buf[:num])
	return num, nil
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (b *Buffer) Next(n int) []byte {
	m := b.Len()
	if n > m {
		n = m
	}

	data := b.buf[b.off : b.off+n]
	b.off += n
	return data
}

// Tidy 整理buffer中的数据, 使得buffer不会无限增长
func (b *Buffer) Tidy() {
	if b.off > 0 {
		var size = len(b.buf) - b.off
		if size > 0 { // 如果size == 0 则不需要重新copy一遍
			copy(b.buf, b.buf[b.off:])
		}

		b.buf = b.buf[:size]
		b.off = 0
	}
}

//func (b *Buffer) MakeCheckpoint() {
//	b.checkpointBuffer, b.checkpointOffset = b.buf, b.off
//}
//
//func (b *Buffer) RestoreCheckpoint() {
//	b.buf, b.off = b.checkpointBuffer, b.checkpointOffset
//}

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}
