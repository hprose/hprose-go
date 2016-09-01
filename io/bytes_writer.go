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
 * io/bytes_writer.go                                     *
 *                                                        *
 * bytes writer for Go.                                   *
 *                                                        *
 * LastModified: Sep 1, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

// BytesWriter implements the io.Writer and io.ByteWriter interfaces by writing
// to a byte slice
type BytesWriter struct {
	buf []byte
}

// Len return the number of bytes of this writer.
func (w *BytesWriter) Len() int {
	return len(w.buf)
}

// Bytes returns the byte slice of this writer.
func (w *BytesWriter) Bytes() []byte {
	return w.buf
}

// String returns the contents of this writer as a string.
// If the ByteWriter is a nil pointer, it returns "<nil>".
func (w *BytesWriter) String() string {
	if w == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(w.buf)
}

// Clear the byte slice of this writer.
func (w *BytesWriter) Clear() {
	w.buf = w.buf[:0]
}

func (w *BytesWriter) grow(n int) int {
	p := len(w.buf)
	c := cap(w.buf)
	l := p + n
	if l > c {
		var buf []byte
		if w.buf == nil && n <= maxSize {
			buf = BytesPool.Get(n)
		} else {
			if l <= maxSize {
				buf = BytesPool.Get(l)
			} else {
				buf = make([]byte, c*2+n)
			}
			copy(buf, w.buf)
			BytesPool.Put(w.buf)
		}
		w.buf = buf
	}
	w.buf = w.buf[:l]
	return p
}

// Grow the the byte slice capacity of this writer.
func (w *BytesWriter) Grow(n int) {
	if n < 0 {
		panic("BytesWriter: negative count")
	}
	p := w.grow(n)
	w.buf = w.buf[0:p]
}

// WriteByte c to the byte slice of this writer.
func (w *BytesWriter) WriteByte(c byte) error {
	w.writeByte(c)
	return nil
}

// Write the contents of b to the byte slice of this writer.
func (w *BytesWriter) Write(b []byte) (int, error) {
	return w.write(b), nil
}

// Close the writer and put the buf to bytes pool
func (w *BytesWriter) Close() {
	BytesPool.Put(w.buf)
}

func (w *BytesWriter) writeByte(c byte) {
	p := w.grow(1)
	w.buf[p] = c
}

func (w *BytesWriter) write(b []byte) int {
	p := w.grow(len(b))
	return copy(w.buf[p:], b)
}

func (w *BytesWriter) writeString(s string) int {
	p := w.grow(len(s))
	return copy(w.buf[p:], s)
}
