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
 * io/reader.go                                           *
 *                                                        *
 * hprose reader for Go.                                  *
 *                                                        *
 * LastModified: Sep 6, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"errors"
)

// Reader is a fine-grained operation struct for Hprose unserialization
// when JSONCompatible is true, the Map data will unserialize to map[string]interface as the default type
type Reader struct {
	RawReader
	Simple         bool
	structRef      []interface{}
	fieldsRef      [][]string
	ref            []interface{}
	JSONCompatible bool
}

// NewReader is the constructor for Hprose Reader
func NewReader(buf []byte, simple bool) (reader *Reader) {
	reader = new(Reader)
	reader.buf = buf
	reader.Simple = simple
	return
}

// CheckTag the next byte in reader is the expected tag or not
func (r *Reader) CheckTag(expectTag byte) (tag byte) {
	tag = r.readByte()
	if tag != expectTag {
		unexpectedTag(tag, []byte{expectTag})
	}
	return
}

// CheckTags the next byte in reader in the expected tags
func (r *Reader) CheckTags(expectTags []byte) (tag byte) {
	tag = r.readByte()
	if bytes.IndexByte(expectTags, tag) == -1 {
		unexpectedTag(tag, expectTags)
	}
	return
}

// ReadBool from the reader
func (r *Reader) ReadBool() bool {
	tag := r.readByte()
	decoder := boolDecoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "bool")
	return false
}

// ReadStringWithoutTag from the reader
func (r *Reader) ReadStringWithoutTag() (str string) {
	str = readString(&r.ByteReader)
	if !r.Simple {
		setReaderRef(r, str)
	}
	return str
}

// ReadString from the reader
func (r *Reader) ReadString() (str string) {
	return ""
}

// ReadRef from the reader
func (r *Reader) ReadRef() interface{} {
	if r.Simple {
		panic(errors.New("reference unserialization can't support in simple mode"))
	}
	return readRef(r, readInt(&r.ByteReader))
}

// private function

func setReaderRef(r *Reader, o interface{}) {
	if r.ref == nil {
		r.ref = make([]interface{}, 0, 64)
	}
	r.ref = append(r.ref, o)
}

func readRef(r *Reader, i int) interface{} {
	return r.ref[i]
}

func resetReaderRef(r *Reader) {
	if r.ref != nil {
		r.ref = r.ref[:0]
	}
}

func tagToString(tag byte) string {
	switch tag {
	case '0':
	case '1':
	case '2':
	case '3':
	case '4':
	case '5':
	case '6':
	case '7':
	case '8':
	case '9':
	case TagInteger:
		return "int"
	case TagLong:
		return "big.Int"
	case TagDouble:
		return "float64"
	case TagNull:
		return "nil"
	case TagEmpty:
		return "empty string"
	case TagTrue:
		return "true"
	case TagFalse:
		return "false"
	case TagNaN:
		return "NaN"
	case TagInfinity:
		return "Infinity"
	case TagDate:
		return "time.Time"
	case TagTime:
		return "time.Time"
	case TagBytes:
		return "[]byte"
	case TagUTF8Char:
		return "string"
	case TagString:
		return "string"
	case TagGUID:
		return "GUID"
	case TagList:
		return "slice"
	case TagMap:
		return "map"
	case TagClass:
		return "struct"
	case TagObject:
		return "struct"
	case TagRef:
		return "reference"
	default:
		unexpectedTag(tag, nil)
	}
	return ""
}
func castError(tag byte, descType string) {
	srcType := tagToString(tag)
	panic(errors.New("can't convert " + srcType + " to " + descType))
}
