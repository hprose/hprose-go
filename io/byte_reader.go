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
 * LastModified: Sep 6, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"errors"
	"io"
	"math"
	"math/big"
	"strconv"
)

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

// UnreadByte unreads 1 byte from the current position.
func (r *ByteReader) UnreadByte() error {
	if r.off > 0 {
		r.off--
	}
	return nil
}

func (r *ByteReader) unreadByte() {
	if r.off > 0 {
		r.off--
	}
}

// Unread n bytes from the current position.
func (r *ByteReader) Unread(n int) {
	if r.off >= n {
		r.off -= n
	} else {
		r.off = 0
	}
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

// ReadInt64 from ByteReader
func ReadInt64(r *ByteReader, tag byte) (i int64) {
	i = 0
	b := r.readByte()
	if b == tag {
		return
	}
	neg := false
	switch b {
	case '-':
		neg = true
		fallthrough
	case '+':
		b = r.readByte()
	}
	if neg {
		for b != tag {
			i = i*10 - int64(b-'0')
			b = r.readByte()
		}
	} else {
		for b != tag {
			i = i*10 + int64(b-'0')
			b = r.readByte()
		}
	}
	return
}

// ReadUint64 from ByteReader
func ReadUint64(r *ByteReader, tag byte) (i uint64) {
	return uint64(ReadInt64(r, tag))
}

// ReadLongAsFloat64 from ByteReader
func ReadLongAsFloat64(r *ByteReader) (f float64) {
	f = 0
	b := r.readByte()
	if b == TagSemicolon {
		return
	}
	neg := false
	switch b {
	case '-':
		neg = true
		fallthrough
	case '+':
		b = r.readByte()
	}
	if neg {
		for b != TagSemicolon {
			f = f*10 - float64(b-'0')
			b = r.readByte()
		}
	} else {
		for b != TagSemicolon {
			f = f*10 + float64(b-'0')
			b = r.readByte()
		}
	}
	return
}

func readInt(r *ByteReader) int {
	return int(ReadInt64(r, TagSemicolon))
}

func readLength(r *ByteReader) int {
	return int(ReadInt64(r, TagQuote))
}

func readCount(r *ByteReader) int {
	return int(ReadInt64(r, TagOpenbrace))
}

func readUntil(r *ByteReader, tag byte) (result []byte) {
	result = r.buf[r.off:]
	i := bytes.IndexByte(result, tag)
	if i < 0 {
		r.off = len(r.buf)
		return
	}
	r.off += i + 1
	return result[:i]
}

func readFloat32(r *ByteReader) float32 {
	s := byteString(readUntil(r, TagSemicolon))
	f, e := strconv.ParseFloat(s, 32)
	if e != nil {
		panic(e)
	}
	return float32(f)
}

func readFloat64(r *ByteReader) float64 {
	s := byteString(readUntil(r, TagSemicolon))
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		panic(e)
	}
	return f
}

func readUTF8Slice(r *ByteReader, length int) []byte {
	var empty = []byte{}
	if length == 0 {
		return empty
	}
	p := r.off
	for i := 0; i < length; i++ {
		b := r.buf[r.off]
		switch b >> 4 {
		case 0, 1, 2, 3, 4, 5, 6, 7:
			r.off++
		case 12, 13:
			r.off += 2
		case 14:
			r.off += 3
		case 15:
			if b&8 == 8 {
				panic(errors.New("bad utf-8 encode"))
			}
			r.off += 4
			i++
		default:
			panic(errors.New("bad utf-8 encode"))
		}
	}
	return r.buf[p:r.off]
}

func readUTF8String(r *ByteReader, length int) string {
	return string(readUTF8Slice(r, length))
}

func readString(r *ByteReader) (result string) {
	result = readUTF8String(r, readLength(r))
	r.readByte()
	return
}

func readUTF8Char(r *ByteReader) (result string) {
	return readUTF8String(r, 1)
}

func readBigInt(r *ByteReader) (result *big.Int) {
	result = new(big.Int)
	result.SetString(byteString(readUntil(r, TagSemicolon)), 10)
	return
}

func readBigFloat(r *ByteReader) (result *big.Float) {
	result, _, _ = new(big.Float).Parse(byteString(readUntil(r, TagSemicolon)), 10)
	return
}

func readInf(r *ByteReader) float64 {
	// '+' - '+' == 0 >= 0, return positive infinity
	// '+' - '-' == -2 < 0, return negative infinity
	return math.Inf(int(TagPos - r.readByte()))
}

func readBytes(r *ByteReader) (result []byte) {
	result = make([]byte, readLength(r))
	r.Read(result)
	r.readByte()
	return
}

func read2Digit(r *ByteReader) int {
	n := int(r.readByte() - '0')
	return n*10 + int(r.readByte()-'0')
}

func read4Digit(r *ByteReader) int {
	n := int(r.readByte() - '0')
	n = n*10 + int(r.readByte()-'0')
	n = n*10 + int(r.readByte()-'0')
	return n*10 + int(r.readByte()-'0')
}

func readTime(r *ByteReader) (hour int, min int, sec int, nsec int, tag byte) {
	hour = read2Digit(r)
	min = read2Digit(r)
	sec = read2Digit(r)
	tag = r.readByte()
	if tag == TagPoint {
		nsec = int(r.readByte() - '0')
		nsec = nsec*10 + int(r.readByte()-'0')
		nsec = nsec*10 + int(r.readByte()-'0')
		nsec *= 1000000
		tag = r.readByte()
		if (tag >= '0') && (tag <= '9') {
			nsec += int(tag-'0') * 100000
			nsec += int(r.readByte()-'0') * 10000
			nsec += int(r.readByte()-'0') * 1000
			tag = r.readByte()
			if (tag >= '0') && (tag <= '9') {
				nsec += int(tag-'0') * 100
				nsec += int(r.readByte()-'0') * 10
				nsec += int(r.readByte() - '0')
				tag = r.readByte()
			}
		}
	}
	return
}

func readDate(r *ByteReader) (year int, month int, day int, tag byte) {
	year = read4Digit(r)
	month = read2Digit(r)
	day = read2Digit(r)
	tag = r.readByte()
	return
}
