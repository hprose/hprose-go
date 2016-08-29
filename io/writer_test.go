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
 * io/writer_test.go                                      *
 *                                                        *
 * hprose writer test for Go.                             *
 *                                                        *
 * LastModified: Aug 29, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"container/list"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestSerializeNil(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(nil)
	if b.String() != "n" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeNil(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(nil)
	}
}

func BenchmarkWriteNil(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteNil()
	}
}

func TestSerializeTrue(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(true)
	if b.String() != "t" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeTrue(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(true)
	}
}

func BenchmarkWriteTrue(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteBool(true)
	}
}

func TestSerializeFalse(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(false)
	if b.String() != "f" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeFalse(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(false)
	}
}

func BenchmarkWriteFalse(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteBool(false)
	}
}

func TestSerializeDigit(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(i)
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
}

func TestSerializeInt(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt32-10) + 10
		writer.Serialize(x)
		if b.String() != "i"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt64-math.MaxInt32-1) + math.MaxInt32 + 1
		writer.Serialize(x)
		if b.String() != "l"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
}

func BenchmarkSerializeInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(i)
	}
}

func BenchmarkWriteInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteInt(int64(i))
	}
}

func TestSerializeInt8(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(int8(i))
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
	for i := 10; i <= 127; i++ {
		b.Truncate(0)
		writer.Serialize(int8(i))
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
	for i := -128; i < 0; i++ {
		b.Truncate(0)
		writer.Serialize(int8(i))
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeInt16(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(int16(math.MaxInt16))
	if b.String() != "i"+strconv.Itoa(math.MaxInt16)+";" {
		t.Error(b.String())
	}
}

func TestSerializeInt32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(int32(math.MaxInt32))
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeInt32(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(int32(i))
	}
}

func TestSerializeInt64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(int64(math.MaxInt32))
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(int64(math.MaxInt64))
	if b.String() != "l"+strconv.Itoa(math.MaxInt64)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt32-10) + 10
		writer.Serialize(uint(x))
		if b.String() != "i"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
	for i := 0; i <= 100; i++ {
		b.Truncate(0)
		x := rand.Intn(math.MaxInt64-math.MaxInt32-1) + math.MaxInt32 + 1
		writer.Serialize(uint(x))
		if b.String() != "l"+strconv.Itoa(x)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeUint8(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	for i := 0; i <= 9; i++ {
		b.Truncate(0)
		writer.Serialize(uint8(i))
		if b.String() != strconv.Itoa(i) {
			t.Error(b.String())
		}
	}
	for i := 10; i <= 255; i++ {
		b.Truncate(0)
		writer.Serialize(uint8(i))
		if b.String() != "i"+strconv.Itoa(i)+";" {
			t.Error(b.String())
		}
	}
}

func TestSerializeUint16(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(uint16(math.MaxUint16))
	if b.String() != "i"+strconv.Itoa(math.MaxUint16)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(uint32(math.MaxUint32))
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(uint64(math.MaxUint32))
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(uint64(math.MaxUint64))
	if b.String() != "l"+strconv.FormatUint(math.MaxUint64, 10)+";" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeUint64(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(uint64(i))
	}
}

func BenchmarkWriteUint64(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteUint(uint64(i))
	}
}

func TestSerializeUintptr(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(uintptr(123))
	if b.String() != "i123;" {
		t.Error(b.String())
	}
}

func TestSerializeFloat32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[float32]string{
		float32(math.NaN()):   "N",
		float32(math.Inf(1)):  "I+",
		float32(math.Inf(-1)): "I-",
		float32(3.14159):      "d3.14159;",
	}
	for k, v := range testdata {
		writer.Serialize(k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func BenchmarkSerializeFloat32(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(float32(i))
	}
}

func BenchmarkWriteFloat32(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteFloat(float64(i), 32)
	}
}

func TestSerializeFloat64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[float64]string{
		math.NaN():       "N",
		math.Inf(1):      "I+",
		math.Inf(-1):     "I-",
		3.14159265358979: "d3.14159265358979;",
	}
	for k, v := range testdata {
		writer.Serialize(k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func BenchmarkSerializeFloat64(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(float64(i))
	}
}

func BenchmarkWriteFloat64(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteFloat(float64(i), 64)
	}
}

func TestSerializeComplex64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(complex(float32(100), 0))
	if b.String() != "d100;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(complex(0, float32(100)))
	if b.String() != "a2{d0;d100;}" {
		t.Error(b.String())
	}
}

func TestSerializeComplex128(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(complex(100, 0))
	if b.String() != "d100;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(complex(0, 100))
	if b.String() != "a2{d0;d100;}" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeComplex128(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(complex(float64(i), float64(i)))
	}
}

func BenchmarkWriteComplex128(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.WriteComplex128(complex(float64(i), float64(i)))
	}
}

func TestWriteTuple(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.WriteTuple()
	if b.String() != "a{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.WriteTuple(1, 3.14, true)
	if b.String() != "a3{1d3.14;t}" {
		t.Error(b.String())
	}
}

func TestWriteBytes(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]byte]string{
		&[]byte{'h', 'e', 'l', 'l', 'o'}: `b5"hello"`,
		&[]byte{}:                        `b""`,
	}
	for k, v := range testdata {
		writer.WriteBytes(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestSerializeString(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[string]string{
		"":                            "e",
		"π":                           "uπ",
		"你":                           "u你",
		"你好":                          `s2"你好"`,
		"你好啊,hello!":                  `s10"你好啊,hello!"`,
		"🇨🇳":                          `s4"🇨🇳"`,
		string([]byte{128, 129, 130}): string([]byte{'b', '3', '"', 128, 129, 130, '"'}),
	}
	for k, v := range testdata {
		writer.Serialize(k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteString(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[string]string{
		"":                            "e",
		"π":                           "uπ",
		"你":                           "u你",
		"你好":                          `s2"你好"`,
		"你好啊,hello!":                  `s10"你好啊,hello!"`,
		"🇨🇳":                          `s4"🇨🇳"`,
		string([]byte{128, 129, 130}): string([]byte{'b', '3', '"', 128, 129, 130, '"'}),
	}
	for k, v := range testdata {
		writer.WriteString(k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestSerializeArray(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[interface{}]string{
		&[...]int{1, 2, 3}:                   "a3{123}",
		&[...]float64{1, 2, 3}:               "a3{d1;d2;d3;}",
		&[...]byte{'h', 'e', 'l', 'l', 'o'}:  `b5"hello"`,
		&[...]byte{}:                         `b""`,
		&[...]interface{}{1, 2.0, nil, true}: "a4{1d2;nt}",
		&[...]bool{true, false, true}:        "a3{tft}",
		&[...]int{}:                          "a{}",
		&[...]bool{}:                         "a{}",
		&[...]interface{}{}:                  "a{}",
	}
	for k, v := range testdata {
		writer.Serialize(k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestSerializeSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[interface{}]string{
		&[]int{1, 2, 3}:                   "a3{123}",
		&[]float64{1, 2, 3}:               "a3{d1;d2;d3;}",
		&[]byte{'h', 'e', 'l', 'l', 'o'}:  `b5"hello"`,
		&[]byte{}:                         `b""`,
		&[]interface{}{1, 2.0, nil, true}: "a4{1d2;nt}",
		&[]bool{true, false, true}:        "a3{tft}",
		&[]int{}:                          "a{}",
		&[]bool{}:                         "a{}",
		&[]interface{}{}:                  "a{}",
	}
	for k, v := range testdata {
		writer.Serialize(k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteBoolSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]bool]string{
		&[]bool{true, false, true}: "a3{tft}",
		&[]bool{}:                  "a{}",
	}
	for k, v := range testdata {
		writer.WriteBoolSlice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteIntSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]int]string{
		&[]int{1, 2, 3}: "a3{123}",
		&[]int{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteIntSlice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteInt8Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]int8]string{
		&[]int8{1, 2, 3}: "a3{123}",
		&[]int8{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteInt8Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteInt16Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]int16]string{
		&[]int16{1, 2, 3}: "a3{123}",
		&[]int16{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteInt16Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteInt32Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]int32]string{
		&[]int32{1, 2, 3}: "a3{123}",
		&[]int32{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteInt32Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteInt64Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]int64]string{
		&[]int64{1, 2, 3}: "a3{123}",
		&[]int64{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteInt64Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteUintSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]uint]string{
		&[]uint{1, 2, 3}: "a3{123}",
		&[]uint{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteUintSlice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteUint8Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]uint8]string{
		&[]uint8{1, 2, 3}: "a3{123}",
		&[]uint8{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteUint8Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteUint16Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]uint16]string{
		&[]uint16{1, 2, 3}: "a3{123}",
		&[]uint16{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteUint16Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteUint32Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]uint32]string{
		&[]uint32{1, 2, 3}: "a3{123}",
		&[]uint32{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteUint32Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteUint64Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]uint64]string{
		&[]uint64{1, 2, 3}: "a3{123}",
		&[]uint64{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteUint64Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteUintptrSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]uintptr]string{
		&[]uintptr{1, 2, 3}: "a3{123}",
		&[]uintptr{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteUintptrSlice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteFloat32Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]float32]string{
		&[]float32{1, 2, 3}: "a3{d1;d2;d3;}",
		&[]float32{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteFloat32Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteFloat64Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]float64]string{
		&[]float64{1, 2, 3}: "a3{d1;d2;d3;}",
		&[]float64{}:        "a{}",
	}
	for k, v := range testdata {
		writer.WriteFloat64Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteComplex64Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]complex64]string{
		&[]complex64{complex(0, 0), complex(1, 0), complex(0, 1)}: "a3{d0;d1;a2{d0;d1;}}",
		&[]complex64{}: "a{}",
	}
	for k, v := range testdata {
		writer.WriteComplex64Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteComplex128Slice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]complex128]string{
		&[]complex128{complex(0, 0), complex(1, 0), complex(0, 1)}: "a3{d0;d1;a2{d0;d1;}}",
		&[]complex128{}:                                            "a{}",
	}
	for k, v := range testdata {
		writer.WriteComplex128Slice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteStringSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[]string]string{
		&[]string{"", "π", "hello"}: `a3{euπs5"hello"}`,
		&[]string{}:                 "a{}",
	}
	for k, v := range testdata {
		writer.WriteStringSlice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestWriteBytesSlice(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[*[][]byte]string{
		&[][]byte{[]byte(""), []byte("π"), []byte("hello")}: `a3{b""b2"π"b5"hello"}`,
		&[][]byte{}: "a{}",
	}
	for k, v := range testdata {
		writer.WriteBytesSlice(*k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func BenchmarkSerializeIntArray(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	array := [...]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.Serialize(array)
	}
}

func BenchmarkSerializeIntSlice(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	slice := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.Serialize(slice)
	}
}

func BenchmarkWriteIntSlice(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	slice := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.WriteIntSlice(slice)
	}
}

func BenchmarkSerializeBytes(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	slice := ([]byte)("你好,hello!")
	for i := 0; i < b.N; i++ {
		writer.Serialize(slice)
	}
}

func BenchmarkWriteBytes(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	slice := ([]byte)("你好,hello!")
	for i := 0; i < b.N; i++ {
		writer.WriteBytes(slice)
	}
}

func BenchmarkSerializeString(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	str := "你好,hello!"
	for i := 0; i < b.N; i++ {
		writer.Serialize(str)
	}
}

func BenchmarkWriteString(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	str := "你好,hello!"
	for i := 0; i < b.N; i++ {
		writer.WriteString(str)
	}
}

func TestSerializeBigInt(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(big.NewInt(123))
	if b.String() != "l123;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(*big.NewInt(123))
	if b.String() != "l123;" {
		t.Error(b.String())
	}
}

func BenchmarkWriteBigInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := big.NewInt(123)
	for i := 0; i < b.N; i++ {
		writer.WriteBigInt(x)
	}
}

func BenchmarkSerializeBigInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := big.NewInt(123)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}

func TestSerializeBigRat(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(big.NewRat(123, 1))
	if b.String() != "l123;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(*big.NewRat(123, 2))
	if b.String() != `s5"123/2"` {
		t.Error(b.String())
	}
}

func BenchmarkWriteBigRat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := big.NewRat(123, 2)
	for i := 0; i < b.N; i++ {
		writer.WriteBigRat(x)
	}
}

func BenchmarkSerializeBigRat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := big.NewRat(123, 2)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}

func TestSerializeBigFloat(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	writer.Serialize(big.NewFloat(3.14159265358979))
	if b.String() != "d3.14159265358979;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(*big.NewFloat(3.14159265358979))
	if b.String() != "d3.14159265358979;" {
		t.Error(b.String())
	}
}

func BenchmarkWriteBigFloat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := big.NewFloat(3.14159265358979)
	for i := 0; i < b.N; i++ {
		writer.WriteBigFloat(x)
	}
}

func BenchmarkSerializeBigFloat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := big.NewFloat(3.14159265358979)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}

func TestWriteTime(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[time.Time]string{
		time.Date(1980, 12, 1, 0, 0, 0, 0, time.UTC):              "D19801201Z",
		time.Date(1970, 1, 1, 12, 34, 56, 0, time.UTC):            "T123456Z",
		time.Date(1970, 1, 1, 12, 34, 56, 789000000, time.UTC):    "T123456.789Z",
		time.Date(1970, 1, 1, 12, 34, 56, 789456000, time.UTC):    "T123456.789456Z",
		time.Date(1970, 1, 1, 12, 34, 56, 789456123, time.UTC):    "T123456.789456123Z",
		time.Date(1980, 12, 1, 12, 34, 56, 0, time.UTC):           "D19801201T123456Z",
		time.Date(1980, 12, 1, 12, 34, 56, 789000000, time.UTC):   "D19801201T123456.789Z",
		time.Date(1980, 12, 1, 12, 34, 56, 789456000, time.UTC):   "D19801201T123456.789456Z",
		time.Date(1980, 12, 1, 12, 34, 56, 789456123, time.UTC):   "D19801201T123456.789456123Z",
		time.Date(1980, 12, 1, 0, 0, 0, 0, time.Local):            "D19801201;",
		time.Date(1970, 1, 1, 12, 34, 56, 0, time.Local):          "T123456;",
		time.Date(1980, 12, 1, 12, 34, 56, 0, time.Local):         "D19801201T123456;",
		time.Date(1980, 12, 1, 12, 34, 56, 789456123, time.Local): "D19801201T123456.789456123;",
	}
	for k, v := range testdata {
		writer.WriteTime(&k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func TestSerializeTime(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	testdata := map[time.Time]string{
		time.Date(1980, 12, 1, 0, 0, 0, 0, time.UTC):              "D19801201Z",
		time.Date(1970, 1, 1, 12, 34, 56, 0, time.UTC):            "T123456Z",
		time.Date(1970, 1, 1, 12, 34, 56, 789000000, time.UTC):    "T123456.789Z",
		time.Date(1970, 1, 1, 12, 34, 56, 789456000, time.UTC):    "T123456.789456Z",
		time.Date(1970, 1, 1, 12, 34, 56, 789456123, time.UTC):    "T123456.789456123Z",
		time.Date(1980, 12, 1, 12, 34, 56, 0, time.UTC):           "D19801201T123456Z",
		time.Date(1980, 12, 1, 12, 34, 56, 789000000, time.UTC):   "D19801201T123456.789Z",
		time.Date(1980, 12, 1, 12, 34, 56, 789456000, time.UTC):   "D19801201T123456.789456Z",
		time.Date(1980, 12, 1, 12, 34, 56, 789456123, time.UTC):   "D19801201T123456.789456123Z",
		time.Date(1980, 12, 1, 0, 0, 0, 0, time.Local):            "D19801201;",
		time.Date(1970, 1, 1, 12, 34, 56, 0, time.Local):          "T123456;",
		time.Date(1980, 12, 1, 12, 34, 56, 0, time.Local):         "D19801201T123456;",
		time.Date(1980, 12, 1, 12, 34, 56, 789456123, time.Local): "D19801201T123456.789456123;",
	}
	for k, v := range testdata {
		writer.Serialize(&k)
		if b.String() != v {
			t.Error(b.String())
		}
		b.Truncate(0)
	}
}

func BenchmarkWriteTime(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := time.Date(1980, 12, 1, 12, 34, 56, 789456123, time.UTC)
	for i := 0; i < b.N; i++ {
		writer.WriteTime(&x)
	}
}

func BenchmarkSerializeTime(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	x := time.Date(1980, 12, 1, 12, 34, 56, 789456123, time.UTC)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}

func TestSerializeList(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	lst := list.New()
	writer.Serialize(lst)
	if b.String() != "a{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	lst.PushBack(1)
	lst.PushBack("hello")
	lst.PushBack(nil)
	lst.PushBack(3.14159)
	writer.Serialize(lst)
	if b.String() != `a4{1s5"hello"nd3.14159;}` {
		t.Error(b.String())
	}
}

func TestWriteList(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	lst := list.New()
	writer.WriteList(lst)
	if b.String() != "a{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	lst.PushBack(1)
	lst.PushBack("hello")
	lst.PushBack(nil)
	lst.PushBack(3.14159)
	writer.WriteList(lst)
	if b.String() != `a4{1s5"hello"nd3.14159;}` {
		t.Error(b.String())
	}
}

func BenchmarkSerializeList(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	lst := list.New()
	lst.PushBack(1)
	lst.PushBack("hello")
	lst.PushBack(nil)
	lst.PushBack(3.14159)
	for i := 0; i < b.N; i++ {
		writer.Serialize(lst)
	}
}

func BenchmarkWriteList(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	lst := list.New()
	lst.PushBack(1)
	lst.PushBack("hello")
	lst.PushBack(nil)
	lst.PushBack(3.14159)
	for i := 0; i < b.N; i++ {
		writer.WriteList(lst)
	}
}

func TestWriterMap(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, true)
	m := make(map[string]interface{})
	writer.Serialize(m)
	if b.String() != "m{}" {
		t.Error(b.String())
	}
	b.Truncate(0)
	m["name"] = "Tom"
	m["age"] = 36
	m["male"] = true
	writer.Serialize(m)
	s := b.String()
	s1 := `m3{s4"name"s3"Tom"s3"age"i36;s4"male"t}`
	s2 := `m3{s3"age"i36;s4"male"ts4"name"s3"Tom"}`
	s3 := `m3{s3"age"i36;s4"name"s3"Tom"s4"male"t}`
	s4 := `m3{s4"name"s3"Tom"s4"male"ts3"age"i36;}`
	s5 := `m3{s4"male"ts3"age"i36;s4"name"s3"Tom"}`
	s6 := `m3{s4"male"ts4"name"s3"Tom"s3"age"i36;}`
	if s != s1 && s != s2 && s != s3 && s != s4 && s != s5 && s != s6 {
		t.Error(b.String())
	}
}

func TestWriterMapRef(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	m := make(map[string]interface{})
	writer.Serialize(m)
	if b.String() != "m{}" {
		t.Error(b.String())
	}
	writer.Reset()
	b.Truncate(0)
	m["name"] = "Tom"
	m["age"] = 36
	m["male"] = true
	writer.Serialize(&m)
	writer.Serialize(&m)
	s := b.String()
	s1 := `m3{s4"name"s3"Tom"s3"age"i36;s4"male"t}r0;`
	s2 := `m3{s3"age"i36;s4"male"ts4"name"s3"Tom"}r0;`
	s3 := `m3{s3"age"i36;s4"name"s3"Tom"s4"male"t}r0;`
	s4 := `m3{s4"name"s3"Tom"s4"male"ts3"age"i36;}r0;`
	s5 := `m3{s4"male"ts3"age"i36;s4"name"s3"Tom"}r0;`
	s6 := `m3{s4"male"ts4"name"s3"Tom"s3"age"i36;}r0;`
	if s != s1 && s != s2 && s != s3 && s != s4 && s != s5 && s != s6 {
		t.Error(b.String())
	}
}

func BenchmarkSerializeStringKeyMap(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	m := make(map[string]interface{})
	m["name"] = "Tom"
	m["age"] = 36
	m["male"] = true
	for i := 0; i < b.N; i++ {
		writer.Serialize(m)
	}
}

func BenchmarkSerializeEmptyMap(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	m := make(map[string]interface{})
	for i := 0; i < b.N; i++ {
		writer.Serialize(m)
	}
}

func BenchmarkSerializeInterfaceKeyMap(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	m := make(map[interface{}]interface{})
	m["name"] = "Tom"
	m["age"] = 36
	m["male"] = true
	for i := 0; i < b.N; i++ {
		writer.Serialize(&m)
	}
}

func TestSerializeStruct(t *testing.T) {
	type TestStruct struct {
		ID int `hprose:"id"`
	}
	type TestStruct1 struct {
		TestStruct
		Name string
		Age  *int
	}
	type TestStruct2 struct {
		OOXX bool `hprose:"ooxx"`
		*TestStruct2
		TestStruct1
		Test     TestStruct
		birthday time.Time
	}
	st := TestStruct2{}
	st.TestStruct2 = &st
	st.ID = 100
	st.Name = "Tom"
	age := 18
	st.Age = &age
	st.OOXX = false
	st.Test.ID = 200
	Register(reflect.TypeOf((*TestStruct)(nil)), "Test", "hprose")
	Register(reflect.TypeOf((*TestStruct1)(nil)), "Test1", "hprose")
	Register(reflect.TypeOf((*TestStruct2)(nil)), "Test2", "hprose")
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	writer.Serialize(st)
	s := `c5"Test2"6{s4"ooxx"s11"testStruct2"s2"id"s4"name"s3"age"s4"test"}o0{fo0{fr7;i100;s3"Tom"i18;c4"Test"1{s2"id"}o1{i200;}}i100;s3"Tom"i18;o1{i200;}}`
	if buf.String() != s {
		t.Error(buf.String())
	}
}

func BenchmarkSerializeStruct(b *testing.B) {
	type TestStruct struct {
		ID int `hprose:"id"`
	}
	type TestStruct1 struct {
		TestStruct
		Name string
		Age  *int
	}
	type TestStruct2 struct {
		OOXX bool `hprose:"ooxx"`
		//*TestStruct2
		TestStruct1
		Test     TestStruct
		birthday time.Time
	}
	st := TestStruct2{}
	//st.TestStruct2 = &st
	st.ID = 100
	st.Name = "Tom"
	age := 18
	st.Age = &age
	st.OOXX = false
	st.Test.ID = 200
	Register(reflect.TypeOf((*TestStruct)(nil)), "Test", "hprose")
	Register(reflect.TypeOf((*TestStruct1)(nil)), "Test1", "hprose")
	Register(reflect.TypeOf((*TestStruct2)(nil)), "Test2", "hprose")
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, true)
	for i := 0; i < b.N; i++ {
		writer.Serialize(st)
	}
}
