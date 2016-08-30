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
 * io/raw_reader.go                                       *
 *                                                        *
 * hprose raw reader for Go.                              *
 *                                                        *
 * LastModified: Aug 29, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"errors"
	"io"
)

// RawReader is the hprose raw reader
type RawReader struct {
	ByteReader
}

// NewRawReader is a constructor for RawReader
func NewRawReader(buf []byte) (reader *RawReader) {
	reader = new(RawReader)
	reader.buf = buf
	reader.len = len(buf)
	return
}

// ReadRaw from stream
func (r *RawReader) ReadRaw() (raw []byte, err error) {
	buffer := new(bytes.Buffer)
	err = r.ReadRawTo(buffer)
	raw = buffer.Bytes()
	return
}

// ReadRawTo buffer from stream
func (r *RawReader) ReadRawTo(buffer *bytes.Buffer) error {
	if r.off >= r.len {
		return io.EOF
	}
	return r.readRaw(buffer, r.readByte())
}

func (r *RawReader) readRaw(buffer *bytes.Buffer, tag byte) (err error) {
	buffer.WriteByte(tag)
	switch tag {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		TagNull, TagEmpty, TagTrue, TagFalse, TagNaN:
	case TagInfinity:
		if r.off >= r.len {
			return io.EOF
		}
		buffer.WriteByte(r.readByte())
	case TagInteger, TagLong, TagDouble, TagRef:
		err = r.readNumberRaw(buffer)
	case TagDate, TagTime:
		err = r.readDateTimeRaw(buffer)
	case TagUTF8Char:
		err = r.readUTF8CharRaw(buffer)
	case TagBytes:
		err = r.readBytesRaw(buffer)
	case TagString:
		err = r.readStringRaw(buffer)
	case TagGUID:
		err = r.readGUIDRaw(buffer)
	case TagList, TagMap, TagObject:
		err = r.readComplexRaw(buffer)
	case TagClass:
		if err = r.readComplexRaw(buffer); err == nil {
			err = r.ReadRawTo(buffer)
		}
	case TagError:
		err = r.ReadRawTo(buffer)
	default:
		err = unexpectedTag(tag, nil)
	}
	return
}

func (r *RawReader) readNumberRaw(buffer *bytes.Buffer) error {
	for r.off < r.len {
		tag := r.readByte()
		buffer.WriteByte(tag)
		if tag == TagSemicolon {
			return nil
		}
	}
	return io.EOF
}

func (r *RawReader) readDateTimeRaw(buffer *bytes.Buffer) error {
	for r.off < r.len {
		tag := r.readByte()
		buffer.WriteByte(tag)
		if tag == TagSemicolon || tag == TagUTC {
			return nil
		}
	}
	return io.EOF
}

func (r *RawReader) readUTF8CharRaw(buffer *bytes.Buffer) (err error) {
	var bytes []byte
	if bytes, err = r.readUTF8Slice(1); err == nil {
		buffer.Write(bytes)
	}
	return
}

func (r *RawReader) readBytesRaw(buffer *bytes.Buffer) (err error) {
	count := 0
	tag := byte('0')
	for r.off < r.len {
		count *= 10
		count += int(tag - '0')
		tag = r.readByte()
		buffer.WriteByte(tag)
		if tag == TagQuote {
			count++
			b := r.Next(count)
			if len(b) < count {
				err = io.EOF
			}
			buffer.Write(b)
			return
		}
	}
	return io.EOF
}

func (r *RawReader) readStringRaw(buffer *bytes.Buffer) (err error) {
	count := 0
	tag := byte('0')
	for r.off < r.len {
		count *= 10
		count += int(tag - '0')
		tag = r.readByte()
		buffer.WriteByte(tag)
		if tag == TagQuote {
			var bytes []byte
			if bytes, err = r.readUTF8Slice(count + 1); err == nil {
				buffer.Write(bytes)
			}
			return
		}
	}
	return io.EOF
}

func (r *RawReader) readGUIDRaw(buffer *bytes.Buffer) (err error) {
	guid := r.Next(38)
	if len(guid) < 38 {
		err = io.EOF
	}
	buffer.Write(guid)
	return err
}

func (r *RawReader) readComplexRaw(buffer *bytes.Buffer) (err error) {
	var tag byte
	for r.off < r.len && tag != TagOpenbrace {
		tag = r.readByte()
		buffer.WriteByte(tag)
	}
	if r.off >= r.len {
		return io.EOF
	}
	tag = r.readByte()
	for err == nil && tag != TagClosebrace {
		if err = r.readRaw(buffer, tag); err == nil {
			tag, err = r.ReadByte()
		}
	}
	if err == nil {
		buffer.WriteByte(tag)
	}
	return err
}

func (r *RawReader) readUTF8Slice(length int) ([]byte, error) {
	var empty = []byte{}
	if length == 0 {
		return empty, nil
	}
	p := r.off
	for i := 0; i < length; i++ {
		if r.off >= r.len {
			return nil, io.EOF
		}
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
				return empty, errors.New("bad utf-8 encode")
			}
			r.off += 4
			i++
		default:
			return empty, errors.New("bad utf-8 encode")
		}
	}
	return r.buf[p:r.off], nil
}

func (r *RawReader) readUTF8String(length int) (string, error) {
	buf, err := r.readUTF8Slice(length)
	return string(buf), err
}

// private functions

func unexpectedTag(tag byte, expectTags []byte) error {
	if t := string([]byte{tag}); expectTags == nil {
		return errors.New("Unexpected serialize tag '" + t + "' in stream")
	} else if bytes.IndexByte(expectTags, tag) < 0 {
		return errors.New("Tag '" + string(expectTags) + "' expected, but '" + t + "' found in stream")
	}
	return nil
}
