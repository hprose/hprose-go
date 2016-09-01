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
 * io/util_test.go                                        *
 *                                                        *
 * util test for Go.                                      *
 *                                                        *
 * LastModified: Aug 24, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"math"
	"reflect"
	"strconv"
	"testing"
)

func BenchmarkGetIntBytes(b *testing.B) {
	buf := make([]byte, 20)
	for i := 0; i < b.N; i++ {
		getIntBytes(buf, int64(i))
		getIntBytes(buf, int64(-i))
		getIntBytes(buf, math.MaxInt32-int64(i))
		getIntBytes(buf, math.MinInt32+int64(i))
		getIntBytes(buf, math.MaxInt64-int64(i))
		getIntBytes(buf, math.MinInt64+int64(i))
	}
}

func BenchmarkGetUintBytes(b *testing.B) {
	buf := make([]byte, 20)
	for i := 0; i < b.N; i++ {
		getUintBytes(buf, uint64(i))
		getUintBytes(buf, uint64(-i))
		getUintBytes(buf, math.MaxUint32-uint64(i))
		getUintBytes(buf, math.MaxUint32+uint64(i))
		getUintBytes(buf, math.MaxUint64-uint64(i))
		getUintBytes(buf, math.MaxUint64+uint64(i))
	}
}

func BenchmarkFormatInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.FormatInt(int64(i), 10)
		strconv.FormatInt(int64(-i), 10)
		strconv.FormatInt(math.MaxInt32-int64(i), 10)
		strconv.FormatInt(math.MinInt32+int64(i), 10)
		strconv.FormatInt(math.MaxInt64-int64(i), 10)
		strconv.FormatInt(math.MinInt64+int64(i), 10)
	}
}

func BenchmarkFormatUint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.FormatUint(uint64(i), 10)
		strconv.FormatUint(uint64(-i), 10)
		strconv.FormatUint(math.MaxUint32-uint64(i), 10)
		strconv.FormatUint(math.MaxUint32+uint64(i), 10)
		strconv.FormatUint(math.MaxUint64-uint64(i), 10)
		strconv.FormatUint(math.MaxUint64+uint64(i), 10)
	}
}

func BenchmarkGetIntBytesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i int64
		buf := make([]byte, 20)
		for pb.Next() {
			getIntBytes(buf, i)
			getIntBytes(buf, -i)
			getIntBytes(buf, math.MaxInt32-i)
			getIntBytes(buf, math.MinInt32+i)
			getIntBytes(buf, math.MaxInt64-i)
			getIntBytes(buf, math.MinInt64+i)
			i++
		}
	})
}

func BenchmarkGetUintBytesParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i uint64
		buf := make([]byte, 20)
		for pb.Next() {
			getUintBytes(buf, i)
			getUintBytes(buf, -i)
			getUintBytes(buf, math.MaxUint32-i)
			getUintBytes(buf, math.MaxUint32+i)
			getUintBytes(buf, math.MaxUint64-i)
			getUintBytes(buf, math.MaxUint64+i)
			i++
		}
	})
}

func BenchmarkFormatIntParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i int64
		for pb.Next() {
			strconv.FormatInt(i, 10)
			strconv.FormatInt(-i, 10)
			strconv.FormatInt(math.MaxInt32-i, 10)
			strconv.FormatInt(math.MinInt32+i, 10)
			strconv.FormatInt(math.MaxInt64-i, 10)
			strconv.FormatInt(math.MinInt64+i, 10)
			i++
		}
	})
}

func BenchmarkFormatUintParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var i uint64
		for pb.Next() {
			strconv.FormatUint(i, 10)
			strconv.FormatUint(-i, 10)
			strconv.FormatUint(math.MaxUint32-i, 10)
			strconv.FormatUint(math.MaxUint32+i, 10)
			strconv.FormatUint(math.MaxUint64-i, 10)
			strconv.FormatUint(math.MaxUint64+i, 10)
			i++
		}
	})
}

func TestGetIntBytes(t *testing.T) {
	data := []int64{
		0, 9, 10, 99, 100, 999, 1000, -1000, 10000, -10000,
		123456789, -123456789, math.MaxInt32, math.MinInt32,
		math.MaxInt64, math.MinInt64}
	buf := make([]byte, 20)
	for _, i := range data {
		b := getIntBytes(buf, i)
		if !reflect.DeepEqual(b, []byte(strconv.FormatInt(i, 10))) {
			t.Error("b must be []byte(\"" + strconv.FormatInt(i, 10) + "\")")
		}
	}
}

func TestGetUintBytes(t *testing.T) {
	data := []uint64{
		0, 9, 10, 99, 100, 999, 1000, 10000, 123456789,
		math.MaxInt32, math.MaxUint32, math.MaxInt64, math.MaxUint64}
	buf := make([]byte, 20)
	for _, i := range data {
		b := getUintBytes(buf, i)
		if !reflect.DeepEqual(b, []byte(strconv.FormatUint(i, 10))) {
			t.Error("b must be []byte(\"" + strconv.FormatUint(i, 10) + "\")")
		}
	}
}

func TestUTF16Length(t *testing.T) {
	data := map[string]int{
		"":                            0,
		"π":                           1,
		"你":                           1,
		"你好":                          2,
		"你好啊,hello!":                  10,
		"🇨🇳":                          4,
		string([]byte{128, 129, 130}): -1,
	}
	for k, v := range data {
		if utf16Length(k) != v {
			t.Error("The UTF16Length of \"" + k + "\" must be " + strconv.Itoa(v))
		}
	}
}

func TestPow2roundup(t *testing.T) {
	data := map[int]int{
		0:             0,
		1:             1,
		2:             2,
		3:             4,
		4:             4,
		5:             8,
		7:             8,
		8:             8,
		9:             16,
		15:            16,
		17:            32,
		31:            32,
		33:            64,
		63:            64,
		65:            128,
		127:           128,
		129:           256,
		257:           512,
		513:           1024,
		math.MaxInt16: math.MaxInt16 + 1,
		math.MaxInt32: math.MaxInt32 + 1,
	}
	for k, v := range data {
		if pow2roundup(k) != v {
			t.Error("The pow2roundup of \"" + strconv.Itoa(k) + "\" must be " + strconv.Itoa(v))
		}
	}
}

func TestLog2(t *testing.T) {
	data := map[int]int{
		0:                 0,
		1:                 0,
		2:                 1,
		4:                 2,
		8:                 3,
		16:                4,
		32:                5,
		64:                6,
		128:               7,
		256:               8,
		512:               9,
		1024:              10,
		math.MaxInt16 + 1: 15,
		math.MaxInt32 + 1: 31,
	}
	for k, v := range data {
		if log2(k) != v {
			t.Error("The log2 of \"" + strconv.Itoa(k) + "\" must be " + strconv.Itoa(v) + ", now it is " + strconv.Itoa(log2(k)))
		}
	}
}
