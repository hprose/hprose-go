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
 * LastModified: Sep 12, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"errors"
	"math/big"
	"reflect"
	"time"
)

// Reader is a fine-grained operation struct for Hprose unserialization
// when JSONCompatible is true, the Map data will unserialize to map[string]interface as the default type
type Reader struct {
	RawReader
	Simple         bool
	fieldsRef      [][]*fieldCache
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
	r.ReadValue(v.Elem())
}

// ReadValue from the reader
func (r *Reader) ReadValue(v reflect.Value) {
	tag := r.readByte()
	decoder := valueDecoders[v.Kind()]
	decoder(r, v, tag)
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

// ReadFloat64 from the reader
func (r *Reader) ReadFloat64() float64 {
	tag := r.readByte()
	decoder := float64Decoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "float64")
	return 0
}

// ReadComplex64 from the reader
func (r *Reader) ReadComplex64() complex64 {
	tag := r.readByte()
	decoder := complex64Decoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "complex64")
	return 0
}

// ReadComplex128 from the reader
func (r *Reader) ReadComplex128() complex128 {
	tag := r.readByte()
	decoder := complex128Decoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "complex128")
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
	tag := r.readByte()
	decoder := stringDecoders[tag]
	if decoder != nil {
		return decoder(r)
	}
	castError(tag, "string")
	return ""
}

// ReadBytesWithoutTag from the reader
func (r *Reader) ReadBytesWithoutTag() (b []byte) {
	l := readLength(&r.ByteReader)
	b = make([]byte, l)
	if _, err := r.Read(b); err != nil {
		panic(err)
	}
	r.readByte()
	if !r.Simple {
		setReaderRef(r, b)
	}
	return
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

// ReadBigIntWithoutTag from the reader
func (r *Reader) ReadBigIntWithoutTag() *big.Int {
	b := readUntil(&r.ByteReader, TagSemicolon)
	i, _ := new(big.Int).SetString(byteString(b), 10)
	return i
}

// ReadSliceWithoutTag from the reader
func (r *Reader) ReadSliceWithoutTag() []reflect.Value {
	l := r.ReadCount()
	v := make([]reflect.Value, l, l+1)
	if !r.Simple {
		setReaderRef(r, v)
	}
	for i := 0; i < l; i++ {
		r.ReadValue(v[i])
	}
	r.readByte()
	return v
}

// ReadCount of array, slice, map or struct field
func (r *Reader) ReadCount() int {
	return int(ReadInt64(&r.ByteReader, TagOpenbrace))
}

// Reset the reference counter
func (r *Reader) Reset() {
	if r.fieldsRef != nil {
		r.fieldsRef = r.fieldsRef[:0]
	}
	if r.Simple {
		return
	}
	if r.ref != nil {
		r.ref = r.ref[:0]
	}
}

// private methods & functions

func (r *Reader) readRef() interface{} {
	if r.Simple {
		panic(errors.New("reference unserialization can't support in simple mode"))
	}
	return readRef(r, readInt(&r.ByteReader))
}

func setReaderRef(r *Reader, o interface{}) {
	r.ref = append(r.ref, o)
}

func readRef(r *Reader, i int) interface{} {
	return r.ref[i]
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
