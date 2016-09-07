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
 * LastModified: Sep 7, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"errors"
	"reflect"
	"time"
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

// Unserialize a data from the reader
func (r *Reader) Unserialize(p interface{}) {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr {
		panic(errors.New("Unserialize: argument p must be a pointer"))
	}
	e := v.Elem()
	decoder := valueDecoders[e.Kind()]
	if decoder != nil {
		decoder(r, e)
	}
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

// ReadIntWithoutTag from the reader
func (r *Reader) ReadIntWithoutTag() int {
	return readInt(&r.ByteReader)
}

// ReadInt from the reader
func (r *Reader) ReadInt() int64 {
	tag := r.readByte()
	decoder := intDecoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "int64")
	return 0
}

// ReadUint from the reader
func (r *Reader) ReadUint() uint64 {
	tag := r.readByte()
	decoder := uintDecoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "uint64")
	return 0
}

// ReadFloat32 from the reader
func (r *Reader) ReadFloat32() float32 {
	tag := r.readByte()
	decoder := float32Decoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "float32")
	return 0
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

// ReadDateTimeWithoutTag from the reader
func (r *Reader) ReadDateTimeWithoutTag() (dt time.Time) {
	year, month, day, tag := readDate(&r.ByteReader)
	var hour, min, sec, nsec int
	if tag == TagTime {
		hour, min, sec, nsec, tag = readTime(&r.ByteReader)
	}
	var loc *time.Location
	if tag == TagUTC {
		loc = time.UTC
	} else {
		loc = time.Local
	}
	dt = time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc)
	if !r.Simple {
		setReaderRef(r, &dt)
	}
	return
}

// ReadTimeWithoutTag from the reader
func (r *Reader) ReadTimeWithoutTag() (t time.Time) {
	hour, min, sec, nsec, tag := readTime(&r.ByteReader)
	var loc *time.Location
	if tag == TagUTC {
		loc = time.UTC
	} else {
		loc = time.Local
	}
	t = time.Date(1970, 1, 1, hour, min, sec, nsec, loc)
	if !r.Simple {
		setReaderRef(r, &t)
	}
	return
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

var tagStringMap = map[byte]string{
	'0':         "int",
	'1':         "int",
	'2':         "int",
	'3':         "int",
	'4':         "int",
	'5':         "int",
	'6':         "int",
	'7':         "int",
	'8':         "int",
	'9':         "int",
	TagInteger:  "int",
	TagLong:     "big.Int",
	TagDouble:   "float64",
	TagNull:     "nil",
	TagEmpty:    "empty string",
	TagTrue:     "true",
	TagFalse:    "false",
	TagNaN:      "NaN",
	TagInfinity: "Infinity",
	TagDate:     "time.Time",
	TagTime:     "time.Time",
	TagBytes:    "[]byte",
	TagUTF8Char: "string",
	TagString:   "string",
	TagGUID:     "GUID",
	TagList:     "slice",
	TagMap:      "map",
	TagClass:    "struct",
	TagObject:   "struct",
	TagRef:      "reference",
}

func tagToString(tag byte) (str string) {
	str = tagStringMap[tag]
	if str == "" {
		unexpectedTag(tag, nil)
	}
	return
}

func castError(tag byte, descType string) {
	srcType := tagToString(tag)
	panic(errors.New("can't convert " + srcType + " to " + descType))
}
