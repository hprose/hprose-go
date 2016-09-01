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
 * LastModified: Sep 1, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "io"

// ByteReader implements the io.Reader and io.ByteReader interfaces by reading
// from a byte slice
type ByteReader struct {
	buf []byte
	off int
}

// NewBytesReader is a constructor for ByteReader
func NewBytesReader(buf []byte) (reader *ByteReader) {
	reader = new(ByteReader)
	reader.buf = buf
	return
}

// ReadByte reads and returns a single byte. If no byte is available,
// it returns error io.EOF.
func (r *ByteReader) ReadByte() (byte, error) {
	if r.off >= len(r.buf) {
		return 0, io.EOF
	}
	return r.readByte(), nil
}

func (r *ByteReader) readByte() (b byte) {
	b = r.buf[r.off]
	r.off++
	return
}

// Read reads the next len(b) bytes from the buffer or until the buffer is
// drained. The return value n is the number of bytes read. If the buffer has
// no data, err is io.EOF (unless len(b) is zero); otherwise it is nil.
func (r *ByteReader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	if r.off >= len(r.buf) {
		return 0, io.EOF
	}
	n = copy(b, r.buf[r.off:])
	r.off += n
	return
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by Read.
// If there are fewer than n bytes, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (r *ByteReader) Next(n int) (data []byte) {
	p := r.off + n
	if p > len(r.buf) {
		p = len(r.buf)
	}
	data = r.buf[r.off:p]
	r.off = p
	return
}
