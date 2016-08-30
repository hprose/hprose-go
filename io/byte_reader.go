/**********************************************************\
|                                                          |
|                          hprose                          |
|                                                          |
| Official WebSite: http://www.hprose.com/                 |
|                   http://www.hprose.org/                 |
|                                                          |
\**********************************************************/
/**********************************************************\
 *                                                        *
 * io/byte_reader.go                                      *
 *                                                        *
 * byte reader for Go.                                    *
 *                                                        *
 * LastModified: Aug 29, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "io"

// ByteReader implements the io.Reader and io.ByteReader interfaces by reading
// from a byte slice
type ByteReader struct {
	buffer []byte
	offset int
	length int
}

// NewByteReader is a constructor for ByteReader
func NewByteReader(buf []byte) (reader *ByteReader) {
	reader = new(ByteReader)
	reader.buffer = buf
	reader.length = len(buf)
	return
}

// ReadByte reads and returns a single byte. If no byte is available,
// it returns error io.EOF.
func (r *ByteReader) ReadByte() (byte, error) {
	if r.offset >= r.length {
		return 0, io.EOF
	}
	return r.readByte(), nil
}

func (r *ByteReader) readByte() (b byte) {
	b = r.buffer[r.offset]
	r.offset++
	return
}

// Read reads the next len(b) bytes from the buffer or until the buffer is
// drained. The return value n is the number of bytes read. If the buffer has
// no data, err is io.EOF (unless len(b) is zero); otherwise it is nil.
func (r *ByteReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.offset >= r.length {
		return 0, io.EOF
	}
	n = copy(b, r.buffer[r.offset:])
	r.offset += n
	return
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (r *ByteReader) Next(n int) (data []byte) {
	p := r.offset + n
	if p > r.length {
		p = r.length
	}
	data = r.buffer[r.offset:p]
	r.offset = p
	return
}
