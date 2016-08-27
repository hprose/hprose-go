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
 * LastModified: Aug 27, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "unsafe"

type mapBodyEncoder func(*Writer, unsafe.Pointer)

var mapBodyEncoders = map[uintptr]mapBodyEncoder{
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

func stringStringMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[string]string)(ptr)
	for k, v := range m {
		writer.WriteString(k)
		writer.WriteString(v)
	}
}

func stringInterfaceMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[string]interface{})(ptr)
	for k, v := range m {
		writer.WriteString(k)
		writer.Serialize(v)
	}
}

func stringIntMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[string]int)(ptr)
	for k, v := range m {
		writer.WriteString(k)
		writer.WriteInt(int64(v))
	}
}

func intIntMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[int]int)(ptr)
	for k, v := range m {
		writer.WriteInt(int64(k))
		writer.WriteInt(int64(v))
	}
}

func intStringMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[int]string)(ptr)
	for k, v := range m {
		writer.WriteInt(int64(k))
		writer.WriteString(v)
	}
}

func intInterfaceMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[int]interface{})(ptr)
	for k, v := range m {
		writer.WriteInt(int64(k))
		writer.Serialize(v)
	}
}

func interfaceInterfaceMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[interface{}]interface{})(ptr)
	for k, v := range m {
		writer.Serialize(k)
		writer.Serialize(v)
	}
}

func interfaceIntMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[interface{}]int)(ptr)
	for k, v := range m {
		writer.Serialize(k)
		writer.WriteInt(int64(v))
	}
}

func interfaceStringMapEncoder(writer *Writer, ptr unsafe.Pointer) {
	m := *(*map[interface{}]string)(ptr)
	for k, v := range m {
		writer.Serialize(k)
		writer.WriteString(v)
	}
}
