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
 * io/map_encoder.go                                      *
 *                                                        *
 * hprose map encoder for Go.                             *
 *                                                        *
 * LastModified: Sep 1, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "unsafe"

var mapBodyEncoders = map[uintptr]func(*Writer, unsafe.Pointer){
	getType((map[string]string)(nil)):           stringStringMapEncoder,
	getType((map[string]interface{})(nil)):      stringInterfaceMapEncoder,
	getType((map[string]int)(nil)):              stringIntMapEncoder,
	getType((map[int]int)(nil)):                 intIntMapEncoder,
	getType((map[int]string)(nil)):              intStringMapEncoder,
	getType((map[int]interface{})(nil)):         intInterfaceMapEncoder,
	getType((map[interface{}]interface{})(nil)): interfaceInterfaceMapEncoder,
	getType((map[interface{}]int)(nil)):         interfaceIntMapEncoder,
	getType((map[interface{}]string)(nil)):      interfaceStringMapEncoder,
}

// RegisterMapEncoder for fast serialize custom map type.
// This function is usually used for code generators.
// This function should be called in package init function.
func RegisterMapEncoder(m interface{}, encoder func(*Writer, unsafe.Pointer)) {
	mapBodyEncoders[getType(m)] = encoder
}

func stringStringMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[string]string)(ptr)
	for k, v := range m {
		w.WriteString(k)
		w.WriteString(v)
	}
}

func stringInterfaceMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[string]interface{})(ptr)
	for k, v := range m {
		w.WriteString(k)
		w.Serialize(v)
	}
}

func stringIntMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[string]int)(ptr)
	for k, v := range m {
		w.WriteString(k)
		w.WriteInt(int64(v))
	}
}

func intIntMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[int]int)(ptr)
	for k, v := range m {
		w.WriteInt(int64(k))
		w.WriteInt(int64(v))
	}
}

func intStringMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[int]string)(ptr)
	for k, v := range m {
		w.WriteInt(int64(k))
		w.WriteString(v)
	}
}

func intInterfaceMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[int]interface{})(ptr)
	for k, v := range m {
		w.WriteInt(int64(k))
		w.Serialize(v)
	}
}

func interfaceInterfaceMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[interface{}]interface{})(ptr)
	for k, v := range m {
		w.Serialize(k)
		w.Serialize(v)
	}
}

func interfaceIntMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[interface{}]int)(ptr)
	for k, v := range m {
		w.Serialize(k)
		w.WriteInt(int64(v))
	}
}

func interfaceStringMapEncoder(w *Writer, ptr unsafe.Pointer) {
	m := *(*map[interface{}]string)(ptr)
	for k, v := range m {
		w.Serialize(k)
		w.WriteString(v)
	}
}
