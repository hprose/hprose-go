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
 * LastModified: Sep 1, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "unsafe"

var sliceBodyEncoders = map[uintptr]func(*Writer, unsafe.Pointer){
	getType(([]bool)(nil)):        boolSliceEncoder,
	getType(([]int)(nil)):         intSliceEncoder,
	getType(([]int8)(nil)):        int8SliceEncoder,
	getType(([]int16)(nil)):       int16SliceEncoder,
	getType(([]int32)(nil)):       int32SliceEncoder,
	getType(([]int64)(nil)):       int64SliceEncoder,
	getType(([]uint)(nil)):        uintSliceEncoder,
	getType(([]uint8)(nil)):       uint8SliceEncoder,
	getType(([]uint16)(nil)):      uint16SliceEncoder,
	getType(([]uint32)(nil)):      uint32SliceEncoder,
	getType(([]uint64)(nil)):      uint64SliceEncoder,
	getType(([]uintptr)(nil)):     uintptrSliceEncoder,
	getType(([]float32)(nil)):     float32SliceEncoder,
	getType(([]float64)(nil)):     float64SliceEncoder,
	getType(([]complex64)(nil)):   complex64SliceEncoder,
	getType(([]complex128)(nil)):  complex128SliceEncoder,
	getType(([]string)(nil)):      stringSliceEncoder,
	getType(([][]byte)(nil)):      bytesSliceEncoder,
	getType(([]interface{})(nil)): interfaceSliceEncoder,
}

// RegisterSliceEncoder for fast serialize custom slice type.
// This function is usually used for code generators.
// This function should be called in package init function.
func RegisterSliceEncoder(m interface{}, encoder func(*Writer, unsafe.Pointer)) {
	sliceBodyEncoders[getType(m)] = encoder
}

func boolSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]bool)(ptr)
	for _, e := range slice {
		w.WriteBool(e)
	}
}

func intSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int)(ptr)
	for _, e := range slice {
		w.WriteInt(int64(e))
	}
}

func int8SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int8)(ptr)
	for _, e := range slice {
		w.WriteInt(int64(e))
	}
}

func int16SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int16)(ptr)
	for _, e := range slice {
		w.WriteInt(int64(e))
	}
}

func int32SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int32)(ptr)
	for _, e := range slice {
		w.WriteInt(int64(e))
	}
}

func int64SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]int64)(ptr)
	for _, e := range slice {
		w.WriteInt(e)
	}
}

func uintSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint)(ptr)
	for _, e := range slice {
		w.WriteUint(uint64(e))
	}
}

func uint8SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint8)(ptr)
	for _, e := range slice {
		w.WriteUint(uint64(e))
	}
}

func uint16SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint16)(ptr)
	for _, e := range slice {
		w.WriteUint(uint64(e))
	}
}

func uint32SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint32)(ptr)
	for _, e := range slice {
		w.WriteUint(uint64(e))
	}
}

func uint64SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uint64)(ptr)
	for _, e := range slice {
		w.WriteUint(e)
	}
}

func uintptrSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]uintptr)(ptr)
	for _, e := range slice {
		w.WriteUint(uint64(e))
	}
}

func float32SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]float32)(ptr)
	for _, e := range slice {
		w.WriteFloat(float64(e), 32)
	}
}

func float64SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]float64)(ptr)
	for _, e := range slice {
		w.WriteFloat(e, 64)
	}
}

func complex64SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]complex64)(ptr)
	for _, e := range slice {
		w.WriteComplex64(e)
	}
}

func complex128SliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]complex128)(ptr)
	for _, e := range slice {
		w.WriteComplex128(e)
	}
}

func stringSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]string)(ptr)
	for _, e := range slice {
		w.WriteString(e)
	}
}

func bytesSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[][]byte)(ptr)
	for _, e := range slice {
		w.WriteBytes(e)
	}
}

func interfaceSliceEncoder(w *Writer, ptr unsafe.Pointer) {
	slice := *(*[]interface{})(ptr)
	for _, e := range slice {
		w.Serialize(e)
	}
}
