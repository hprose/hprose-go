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
 * io/byte_writer.go                                      *
 *                                                        *
 * byte writer for Go.                                    *
 *                                                        *
 * LastModified: Aug 29, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "math"

// ByteWriter implements the io.Writer and io.ByteWriter interfaces by writing
// to a byte slice
type ByteWriter struct {
	buf []byte
}

func pow2roundup(x int) int {
	x--
	x |= x >> 1
	x |= x >> 2
	x |= x >> 4
	x |= x >> 8
	x |= x >> 16
	return x + 1
}

// Len return the number of bytes of this writer.
func (w *ByteWriter) Len() int {
	return len(w.buf)
}

// Bytes returns the byte slice of this writer.
func (w *ByteWriter) Bytes() []byte {
	return w.buf
}

// String returns the contents of this writer as a string.
// If the ByteWriter is a nil pointer, it returns "<nil>".
func (w *ByteWriter) String() string {
	if w == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(w.buf)
}

// Clear the byte slice of this writer.
func (w *ByteWriter) Clear() {
	w.buf = w.buf[:0]
}

func (w *ByteWriter) grow(n int) int {
	p := len(w.buf)
	c := cap(w.buf)
	l := p + n
	if l > c {
		var buf []byte
		if w.buf == nil && n <= 64 {
			buf = make([]byte, 64)
		} else {
			if l < math.MaxInt32 {
				buf = make([]byte, pow2roundup(l))
			} else {
				buf = make([]byte, c*2+n)
			}
			copy(buf, w.buf)
		}
		w.buf = buf
	}
	w.buf = w.buf[:l]
	return p
}

// Grow the the byte slice capacity of this writer.
func (w *ByteWriter) Grow(n int) {
	if n < 0 {
		panic("BytesWriter: negative count")
	}
	p := w.grow(n)
	w.buf = w.buf[0:p]
}

// WriteByte c to the byte slice of this writer.
func (w *ByteWriter) WriteByte(c byte) error {
	w.writeByte(c)
	return nil
}

// Write the contents of b to the byte slice of this writer.
func (w *ByteWriter) Write(b []byte) (int, error) {
	return w.write(b), nil
}

func (w *ByteWriter) writeByte(c byte) {
	p := w.grow(1)
	w.buf[p] = c
}

func (w *ByteWriter) write(b []byte) int {
	p := w.grow(len(b))
	return copy(w.buf[p:], b)
}

func (w *ByteWriter) writeString(s string) int {
	p := w.grow(len(s))
	return copy(w.buf[p:], s)
}
