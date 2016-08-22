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
 * io/slice_encoder.go                                    *
 *                                                        *
 * hprose slice encoder for Go.                           *
 *                                                        *
 * LastModified: Aug 22, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"reflect"
	"unsafe"
)

type sliceBodyEncoder func(*Writer, unsafe.Pointer)

var sliceBodyEncoders []sliceBodyEncoder

func boolSliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]bool)(ptr)
	for _, e := range slice {
		writer.WriteBool(e)
	}
}

func intSliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int)(ptr)
	for _, e := range slice {
		writer.WriteInt(int64(e))
	}
}

func int8SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int8)(ptr)
	for _, e := range slice {
		writer.WriteInt32(int32(e))
	}
}

func int16SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int16)(ptr)
	for _, e := range slice {
		writer.WriteInt32(int32(e))
	}
}

func int32SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int32)(ptr)
	for _, e := range slice {
		writer.WriteInt32(e)
	}
}

func int64SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int64)(ptr)
	for _, e := range slice {
		writer.WriteInt(e)
	}
}

func uintSliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint)(ptr)
	for _, e := range slice {
		writer.WriteUint(uint64(e))
	}
}

func uint8SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint8)(ptr)
	for _, e := range slice {
		writer.WriteInt32(int32(e))
	}
}

func uint16SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint16)(ptr)
	for _, e := range slice {
		writer.WriteInt32(int32(e))
	}
}

func uint32SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint32)(ptr)
	for _, e := range slice {
		writer.WriteUint(uint64(e))
	}
}

func uint64SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint64)(ptr)
	for _, e := range slice {
		writer.WriteUint(e)
	}
}

func uintptrSliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uintptr)(ptr)
	for _, e := range slice {
		writer.WriteUint(uint64(e))
	}
}

func float32SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]float32)(ptr)
	for _, e := range slice {
		writer.WriteFloat(float64(e), 32)
	}
}

func float64SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]float64)(ptr)
	for _, e := range slice {
		writer.WriteFloat(e, 64)
	}
}

func complex64SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]complex64)(ptr)
	for _, e := range slice {
		writer.WriteComplex64(e)
	}
}

func complex128SliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]complex128)(ptr)
	for _, e := range slice {
		writer.WriteComplex128(e)
	}
}

func stringSliceEncoder(writer *Writer, ptr unsafe.Pointer) {
	slice := *(*[]string)(ptr)
	for _, e := range slice {
		writer.WriteString(e)
	}
}

func init() {
	sliceBodyEncoders = []sliceBodyEncoder{
		reflect.Invalid:       nil,
		reflect.Bool:          boolSliceEncoder,
		reflect.Int:           intSliceEncoder,
		reflect.Int8:          int8SliceEncoder,
		reflect.Int16:         int16SliceEncoder,
		reflect.Int32:         int32SliceEncoder,
		reflect.Int64:         int64SliceEncoder,
		reflect.Uint:          uintSliceEncoder,
		reflect.Uint8:         uint8SliceEncoder,
		reflect.Uint16:        uint16SliceEncoder,
		reflect.Uint32:        uint32SliceEncoder,
		reflect.Uint64:        uint64SliceEncoder,
		reflect.Uintptr:       uintptrSliceEncoder,
		reflect.Float32:       float32SliceEncoder,
		reflect.Float64:       float64SliceEncoder,
		reflect.Complex64:     complex64SliceEncoder,
		reflect.Complex128:    complex128SliceEncoder,
		reflect.Array:         nil,
		reflect.Chan:          nil,
		reflect.Func:          nil,
		reflect.Interface:     nil,
		reflect.Map:           nil,
		reflect.Ptr:           nil,
		reflect.Slice:         nil,
		reflect.String:        stringSliceEncoder,
		reflect.Struct:        nil,
		reflect.UnsafePointer: nil,
	}
}
