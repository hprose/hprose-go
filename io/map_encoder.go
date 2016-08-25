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
 * LastModified: Aug 22, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "unsafe"

type mapBodyEncoder func(*Writer, unsafe.Pointer)

func getMapEncoder(mapType uintptr) mapBodyEncoder {
	switch mapType {
	case stringStringMapType:
		return stringStringMapEncoder
	case stringInterfaceMapType:
		return stringInterfaceMapEncoder
	case stringIntMapType:
		return stringIntMapEncoder
	case intIntMapType:
		return intIntMapEncoder
	case intStringMapType:
		return intStringMapEncoder
	case intInterfaceMapType:
		return intInterfaceMapEncoder
	default:
		return nil
	}
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
