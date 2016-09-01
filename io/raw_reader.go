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
 * LastModified: Sep 1, 2016                              *
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
	return
}

// ReadRaw from stream
func (r *RawReader) ReadRaw() (raw []byte, err error) {
	w := new(ByteWriter)
	err = r.ReadRawTo(w)
	raw = w.Bytes()
	return
}

// ReadRawTo buffer from stream
func (r *RawReader) ReadRawTo(w *ByteWriter) error {
	if r.off >= len(r.buf) {
		return io.EOF
	}
	return r.readRaw(w, r.readByte())
}

func (r *RawReader) readRaw(w *ByteWriter, tag byte) (err error) {
	w.writeByte(tag)
	switch tag {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
		TagNull, TagEmpty, TagTrue, TagFalse, TagNaN:
	case TagInfinity:
		if r.off >= len(r.buf) {
			return io.EOF
		}
		w.writeByte(r.readByte())
	case TagInteger, TagLong, TagDouble, TagRef:
		err = r.readNumberRaw(w)
	case TagDate, TagTime:
		err = r.readDateTimeRaw(w)
	case TagUTF8Char:
		err = r.readUTF8CharRaw(w)
	case TagBytes:
		err = r.readBytesRaw(w)
	case TagString:
		err = r.readStringRaw(w)
	case TagGUID:
		err = r.readGUIDRaw(w)
	case TagList, TagMap, TagObject:
		err = r.readComplexRaw(w)
	case TagClass:
		if err = r.readComplexRaw(w); err == nil {
			err = r.ReadRawTo(w)
		}
	case TagError:
		err = r.ReadRawTo(w)
	default:
		err = unexpectedTag(tag, nil)
	}
	return
}

func (r *RawReader) readNumberRaw(w *ByteWriter) error {
	for r.off < len(r.buf) {
		tag := r.readByte()
		w.writeByte(tag)
		if tag == TagSemicolon {
			return nil
		}
	}
	return io.EOF
}

func (r *RawReader) readDateTimeRaw(w *ByteWriter) error {
	for r.off < len(r.buf) {
		tag := r.readByte()
		w.writeByte(tag)
		if tag == TagSemicolon || tag == TagUTC {
			return nil
		}
	}
	return io.EOF
}

func (r *RawReader) readUTF8CharRaw(w *ByteWriter) (err error) {
	var bytes []byte
	if bytes, err = r.readUTF8Slice(1); err == nil {
		w.write(bytes)
	}
	return
}

func (r *RawReader) readBytesRaw(w *ByteWriter) (err error) {
	count := 0
	tag := byte('0')
	for r.off < len(r.buf) {
		count *= 10
		count += int(tag - '0')
		tag = r.readByte()
		w.writeByte(tag)
		if tag == TagQuote {
			count++
			b := r.Next(count)
			if len(b) < count {
				err = io.EOF
			}
			w.write(b)
			return
		}
	}
	return io.EOF
}

func (r *RawReader) readStringRaw(w *ByteWriter) (err error) {
	count := 0
	tag := byte('0')
	for r.off < len(r.buf) {
		count *= 10
		count += int(tag - '0')
		tag = r.readByte()
		w.writeByte(tag)
		if tag == TagQuote {
			var bytes []byte
			if bytes, err = r.readUTF8Slice(count + 1); err == nil {
				w.write(bytes)
			}
			return
		}
	}
	return io.EOF
}

func (r *RawReader) readGUIDRaw(w *ByteWriter) (err error) {
	guid := r.Next(38)
	if len(guid) < 38 {
		err = io.EOF
	}
	w.write(guid)
	return err
}

func (r *RawReader) readComplexRaw(w *ByteWriter) (err error) {
	var tag byte
	for r.off < len(r.buf) && tag != TagOpenbrace {
		tag = r.readByte()
		w.writeByte(tag)
	}
	if r.off >= len(r.buf) {
		return io.EOF
	}
	tag = r.readByte()
	for err == nil && tag != TagClosebrace {
		if err = r.readRaw(w, tag); err == nil {
			tag, err = r.ReadByte()
		}
	}
	if err == nil {
		w.writeByte(tag)
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
		if r.off >= len(r.buf) {
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
