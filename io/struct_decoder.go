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
 * io/struct_decoder.go                                   *
 *                                                        *
 * hprose struct decoder for Go.                          *
 *                                                        *
 * LastModified: Sep 10, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"math/big"
	"reflect"
	"unsafe"
)

func readDigitAsStruct(r *Reader, v reflect.Value, tag byte) {
	typ := (*reflectValue)(unsafe.Pointer(&v)).typ
	switch typ {
	case bigIntType:
		v.Set(reflect.ValueOf(*big.NewInt(int64(tag - '0'))))
	case bigRatType:
		v.Set(reflect.ValueOf(*big.NewRat(int64(tag-'0'), 1)))
	case bigFloatType:
		v.Set(reflect.ValueOf(*big.NewFloat(float64(tag - '0'))))
	}
}

func readMapAsStruct(r *Reader, v reflect.Value, tag byte) {
}

func readStructMeta(r *Reader, v reflect.Value, tag byte) {
}

func readStructData(r *Reader, v reflect.Value, tag byte) {
}

func readRefAsStruct(r *Reader, v reflect.Value, tag byte) {
}

var structDecoders = [256]func(r *Reader, v reflect.Value, tag byte){
	'0': readDigitAsStruct,
	'1': readDigitAsStruct,
	'2': readDigitAsStruct,
	'3': readDigitAsStruct,
	'4': readDigitAsStruct,
	'5': readDigitAsStruct,
	'6': readDigitAsStruct,
	'7': readDigitAsStruct,
	'8': readDigitAsStruct,
	'9': readDigitAsStruct,
	// TagInteger: readInt64,
	// TagLong:    readInt64,
	TagNull:   nilDecoder,
	TagMap:    readMapAsStruct,
	TagClass:  readStructMeta,
	TagObject: readStructData,
	TagRef:    readRefAsStruct,
}

func structDecoder(r *Reader, v reflect.Value, tag byte) {
	decoder := structDecoders[tag]
	if decoder != nil {
		decoder(r, v, tag)
		return
	}
	castError(tag, "struct")
}
