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
 * io/encoder.go                                          *
 *                                                        *
 * hprose encoder for Go.                                 *
 *                                                        *
 * LastModified: Aug 25, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"container/list"
	"math/big"
	"reflect"
	"time"
	"unsafe"
)

type valueEncoder func(writer *Writer, v reflect.Value)

var valueEncoders []valueEncoder

func nilEncoder(writer *Writer, v reflect.Value) {
	writer.WriteNil()
}

func boolEncoder(writer *Writer, v reflect.Value) {
	writer.WriteBool(v.Bool())
}

func intEncoder(writer *Writer, v reflect.Value) {
	writer.WriteInt(v.Int())
}

func uintEncoder(writer *Writer, v reflect.Value) {
	writer.WriteUint(v.Uint())
}

func float32Encoder(writer *Writer, v reflect.Value) {
	writer.WriteFloat(v.Float(), 32)
}

func float64Encoder(writer *Writer, v reflect.Value) {
	writer.WriteFloat(v.Float(), 64)
}

func complex64Encoder(writer *Writer, v reflect.Value) {
	writer.WriteComplex64(complex64(v.Complex()))
}

func complex128Encoder(writer *Writer, v reflect.Value) {
	writer.WriteComplex128(v.Complex())
}

func interfaceEncoder(writer *Writer, v reflect.Value) {
	if v.IsNil() {
		writer.WriteNil()
		return
	}
	e := v.Elem()
	valueEncoders[e.Kind()](writer, e)
}

func arrayEncoder(writer *Writer, v reflect.Value) {
	writer.SetRef(nil)
	writeArray(writer, v)
}

func sliceEncoder(writer *Writer, v reflect.Value) {
	writer.SetRef(nil)
	writeSlice(writer, v)
}

func mapEncoder(writer *Writer, v reflect.Value) {
	ptr := ((*emptyInterface)(unsafe.Pointer(&v)).ptr)
	if !writer.WriteRef(ptr) {
		writer.SetRef(ptr)
		writeMap(writer, v)
	}
}

func stringEncoder(writer *Writer, v reflect.Value) {
	writer.WriteString(v.String())
}

func structEncoder(writer *Writer, v reflect.Value) {
	ptr := ((*emptyInterface)(unsafe.Pointer(&v)).ptr)
	pv := reflect.NewAt(v.Type(), ptr)
	structPtrEncoder(writer, pv, ptr)
}

func arrayPtrEncoder(writer *Writer, v reflect.Value, ptr unsafe.Pointer) {
	if !writer.WriteRef(ptr) {
		writer.SetRef(ptr)
		writeArray(writer, v)
	}
}

func mapPtrEncoder(writer *Writer, v reflect.Value, ptr unsafe.Pointer) {
	if !writer.WriteRef(ptr) {
		writer.SetRef(ptr)
		writeMap(writer, v)
	}
}

func slicePtrEncoder(writer *Writer, v reflect.Value, ptr unsafe.Pointer) {
	if !writer.WriteRef(ptr) {
		writer.SetRef(ptr)
		writeSlice(writer, v)
	}
}

func stringPtrEncoder(writer *Writer, v reflect.Value, ptr unsafe.Pointer) {
	str := v.String()
	length := utf16Length(str)
	switch {
	case length == 0:
		writer.Stream.WriteByte(TagEmpty)
	case length < 0:
		writer.WriteBytes(*(*[]byte)(unsafe.Pointer(&str)))
	case length == 1:
		writer.Stream.WriteByte(TagUTF8Char)
		writer.Stream.WriteString(str)
	default:
		if !writer.WriteRef(ptr) {
			writer.SetRef(ptr)
			writeString(writer, str, length)
		}
	}
}

func structPtrEncoder(writer *Writer, v reflect.Value, ptr unsafe.Pointer) {
	switch *(*uintptr)(unsafe.Pointer(&v)) {
	case bigIntType:
		writer.WriteBigInt((*big.Int)(ptr))
	case bigRatType:
		writer.WriteBigRat((*big.Rat)(ptr))
	case bigFloatType:
		writer.WriteBigFloat((*big.Float)(ptr))
	case timeType:
		if !writer.WriteRef(ptr) {
			writer.SetRef(ptr)
			writeTime(writer, (*time.Time)(ptr))
		}
	case listType:
		if !writer.WriteRef(ptr) {
			writer.SetRef(ptr)
			writeList(writer, (*list.List)(ptr))
		}
	default:
		if !writer.WriteRef(ptr) {
			writer.SetRef(ptr)
			//writeStruct(writer, v)
		}
	}
}

func ptrEncoder(writer *Writer, v reflect.Value) {
	if v.IsNil() {
		writer.WriteNil()
		return
	}
	e := v.Elem()
	kind := e.Kind()
	ptr := ((*emptyInterface)(unsafe.Pointer(&v)).ptr)
	switch kind {
	case reflect.Array:
		arrayPtrEncoder(writer, e, ptr)
	case reflect.Map:
		mapPtrEncoder(writer, e, ptr)
	case reflect.Slice:
		slicePtrEncoder(writer, e, ptr)
	case reflect.String:
		stringPtrEncoder(writer, e, ptr)
	case reflect.Struct:
		structPtrEncoder(writer, v, ptr)
	default:
		valueEncoders[kind](writer, e)
	}
}

func init() {
	valueEncoders = []valueEncoder{
		reflect.Invalid:       nilEncoder,
		reflect.Bool:          boolEncoder,
		reflect.Int:           intEncoder,
		reflect.Int8:          intEncoder,
		reflect.Int16:         intEncoder,
		reflect.Int32:         intEncoder,
		reflect.Int64:         intEncoder,
		reflect.Uint:          uintEncoder,
		reflect.Uint8:         uintEncoder,
		reflect.Uint16:        uintEncoder,
		reflect.Uint32:        uintEncoder,
		reflect.Uint64:        uintEncoder,
		reflect.Uintptr:       uintEncoder,
		reflect.Float32:       float32Encoder,
		reflect.Float64:       float64Encoder,
		reflect.Complex64:     complex64Encoder,
		reflect.Complex128:    complex128Encoder,
		reflect.Array:         arrayEncoder,
		reflect.Chan:          nilEncoder,
		reflect.Func:          nilEncoder,
		reflect.Interface:     interfaceEncoder,
		reflect.Map:           mapEncoder,
		reflect.Ptr:           ptrEncoder,
		reflect.Slice:         sliceEncoder,
		reflect.String:        stringEncoder,
		reflect.Struct:        structEncoder,
		reflect.UnsafePointer: nilEncoder,
	}
}
