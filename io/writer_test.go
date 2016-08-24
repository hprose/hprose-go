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
 * LastModified: Aug 22, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
)

func TestSerializeNil(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(nil)
	if b.String() != "n" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeNil(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(nil)
	}
}

func BenchmarkWriteNil(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteNil()
	}
}

func TestSerializeTrue(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(true)
	if b.String() != "t" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeTrue(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(true)
	}
}

func BenchmarkWriteTrue(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteBool(true)
	}
}

func TestSerializeFalse(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(false)
	if b.String() != "f" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeFalse(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(false)
	}
}

func BenchmarkWriteFalse(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteBool(false)
	}
}

func TestSerializeDigit(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(i)
	}
}

func BenchmarkWriteInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteInt(int64(i))
	}
}

func TestSerializeInt8(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
	writer.Serialize(int16(math.MaxInt16))
	if b.String() != "i"+strconv.Itoa(math.MaxInt16)+";" {
		t.Error(b.String())
	}
}

func TestSerializeInt32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(int32(math.MaxInt32))
	if b.String() != "i"+strconv.Itoa(math.MaxInt32)+";" {
		t.Error(b.String())
	}
}

func BenchmarkSerializeInt32(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(int32(i))
	}
}

func TestSerializeInt64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
	writer.Serialize(uint16(math.MaxUint16))
	if b.String() != "i"+strconv.Itoa(math.MaxUint16)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(uint32(math.MaxUint32))
	if b.String() != "l"+strconv.Itoa(math.MaxUint32)+";" {
		t.Error(b.String())
	}
}

func TestSerializeUint64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(uint64(i))
	}
}

func BenchmarkWriteUint64(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteUint(uint64(i))
	}
}

func TestSerializeUintptr(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(uintptr(123))
	if b.String() != "i123;" {
		t.Error(b.String())
	}
}

func TestSerializeFloat32(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(float32(i))
	}
}

func BenchmarkWriteFloat32(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteFloat(float64(i), 32)
	}
}

func TestSerializeFloat64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(float64(i))
	}
}

func BenchmarkWriteFloat64(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteFloat(float64(i), 64)
	}
}

func TestSerializeComplex64(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.Serialize(complex(float64(i), float64(i)))
	}
}

func BenchmarkWriteComplex128(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	for i := 0; i < b.N; i++ {
		writer.WriteComplex128(complex(float64(i), float64(i)))
	}
}

func TestWriteTuple(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
	testdata := map[*[]byte]string{
		&[]byte{'h', 'e', 'l', 'l', 'o'}: "b5\"hello\"",
		&[]byte{}:                        "b\"\"",
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
	writer := NewWriter(b, false)
	testdata := map[string]string{
		"":                            "e",
		"π":                           "uπ",
		"你":                           "u你",
		"你好":                          "s2\"你好\"",
		"你好啊,hello!":                  "s10\"你好啊,hello!\"",
		"🇨🇳":                          "s4\"🇨🇳\"",
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
	writer := NewWriter(b, false)
	testdata := map[string]string{
		"":                            "e",
		"π":                           "uπ",
		"你":                           "u你",
		"你好":                          "s2\"你好\"",
		"你好啊,hello!":                  "s10\"你好啊,hello!\"",
		"🇨🇳":                          "s4\"🇨🇳\"",
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
	writer := NewWriter(b, false)
	testdata := map[interface{}]string{
		&[...]int{1, 2, 3}:                  "a3{123}",
		&[...]float64{1, 2, 3}:              "a3{d1;d2;d3;}",
		&[...]byte{'h', 'e', 'l', 'l', 'o'}: "b5\"hello\"",
		&[...]byte{}:                        "b\"\"",
		&[...]interface{}{1, 2.0, true}:     "a3{1d2;t}",
		&[...]bool{true, false, true}:       "a3{tft}",
		&[...]int{}:                         "a{}",
		&[...]bool{}:                        "a{}",
		&[...]interface{}{}:                 "a{}",
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
	writer := NewWriter(b, false)
	testdata := map[interface{}]string{
		&[]int{1, 2, 3}:                  "a3{123}",
		&[]float64{1, 2, 3}:              "a3{d1;d2;d3;}",
		&[]byte{'h', 'e', 'l', 'l', 'o'}: "b5\"hello\"",
		&[]byte{}:                        "b\"\"",
		&[]interface{}{1, 2.0, true}:     "a3{1d2;t}",
		&[]bool{true, false, true}:       "a3{tft}",
		&[]int{}:                         "a{}",
		&[]bool{}:                        "a{}",
		&[]interface{}{}:                 "a{}",
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
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
	writer := NewWriter(b, false)
	testdata := map[*[]string]string{
		&[]string{"", "π", "hello"}: "a3{euπs5\"hello\"}",
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
	writer := NewWriter(b, false)
	testdata := map[*[][]byte]string{
		&[][]byte{[]byte(""), []byte("π"), []byte("hello")}: "a3{b\"\"b2\"π\"b5\"hello\"}",
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

func BenchmarkSerializeArray(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	array := [...]int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.Serialize(array)
	}
}

func BenchmarkSerializeSlice(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	slice := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.Serialize(slice)
	}
}

func BenchmarkWriteIntSlice(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	slice := []int{0, 1, 2, 3, 4, 0, 1, 2, 3, 4, 1, 2, 3, 4, 0, 1, 2, 3, 4}
	for i := 0; i < b.N; i++ {
		writer.WriteIntSlice(slice)
	}
}

func BenchmarkSerializeBytes(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	slice := ([]byte)("你好,hello!")
	for i := 0; i < b.N; i++ {
		writer.Serialize(slice)
	}
}

func BenchmarkWriteBytes(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	slice := ([]byte)("你好,hello!")
	for i := 0; i < b.N; i++ {
		writer.WriteBytes(slice)
	}
}

func BenchmarkSerializeString(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	str := "你好,hello!"
	for i := 0; i < b.N; i++ {
		writer.Serialize(str)
	}
}

func BenchmarkWriteString(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	str := "你好,hello!"
	for i := 0; i < b.N; i++ {
		writer.WriteString(str)
	}
}

func TestSerializeBigInt(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	x := big.NewInt(123)
	for i := 0; i < b.N; i++ {
		writer.WriteBigInt(x)
	}
}

func BenchmarkSerializeBigInt(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	x := big.NewInt(123)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}

func TestSerializeBigRat(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
	writer.Serialize(big.NewRat(123, 1))
	if b.String() != "l123;" {
		t.Error(b.String())
	}
	b.Truncate(0)
	writer.Serialize(*big.NewRat(123, 2))
	if b.String() != "s5\"123/2\"" {
		t.Error(b.String())
	}
}

func BenchmarkWriteBigRat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	x := big.NewRat(123, 2)
	for i := 0; i < b.N; i++ {
		writer.WriteBigRat(x)
	}
}

func BenchmarkSerializeBigRat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	x := big.NewRat(123, 2)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}

func TestSerializeBigFloat(t *testing.T) {
	b := new(bytes.Buffer)
	writer := NewWriter(b, false)
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
	writer := NewWriter(buf, false)
	x := big.NewFloat(3.14159265358979)
	for i := 0; i < b.N; i++ {
		writer.WriteBigFloat(x)
	}
}

func BenchmarkSerializeBigFloat(b *testing.B) {
	buf := new(bytes.Buffer)
	writer := NewWriter(buf, false)
	x := big.NewFloat(3.14159265358979)
	for i := 0; i < b.N; i++ {
		writer.Serialize(x)
	}
}
