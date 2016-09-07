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
 * io/reader_test.go                                      *
 *                                                        *
 * hprose Reader Test for Go.                             *
 *                                                        *
 * LastModified: Sep 6, 2016                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"math"
	"testing"
	"time"
)

func TestReadBool(t *testing.T) {
	trueValue := "true"
	data := map[interface{}]bool{
		true:            true,
		false:           false,
		nil:             false,
		"":              false,
		0:               false,
		1:               true,
		9:               true,
		100:             true,
		100000000000000: true,
		0.0:             false,
		"t":             true,
		"f":             false,
		&trueValue:      true,
		&trueValue:      true,
		"false":         false,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&trueValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		b := reader.ReadBool()
		if b != data[k] {
			t.Error(k, data[k], b)
		}
	}
	b := reader.ReadBool()
	if b != true {
		t.Error(trueValue, true, b)
	}
	w.Close()
}

func TestUnserializeBool(t *testing.T) {
	trueValue := "true"
	data := map[interface{}]bool{
		true:            true,
		false:           false,
		nil:             false,
		"":              false,
		0:               false,
		1:               true,
		9:               true,
		100:             true,
		100000000000000: true,
		0.0:             false,
		"t":             true,
		"f":             false,
		&trueValue:      true,
		"false":         false,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&trueValue)
	reader := NewReader(w.Bytes(), false)
	var p bool
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != true {
		t.Error(trueValue, true, p)
	}
	w.Close()
}

func BenchmarkReadBool(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(true)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadBool()
	}
	w.Close()
}

func BenchmarkUnserializeBool(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(true)
	bytes := w.Bytes()
	var p bool
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}

func TestReadInt(t *testing.T) {
	intValue := "1234567"
	u := uint(math.MaxUint64)
	data := map[interface{}]int64{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		-100:                             -100,
		math.MinInt32:                    int64(math.MinInt32),
		math.MaxInt64:                    int64(math.MaxInt64),
		math.MinInt64:                    int64(math.MinInt64),
		u:                                int64(u),
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&intValue:                        1234567,
		time.Unix(123, 456):              123000000456,
		time.Unix(1234567890, 123456789): 1234567890123456789,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&intValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		i := reader.ReadInt()
		if i != data[k] {
			t.Error(k, data[k], i)
		}
	}
	i := reader.ReadInt()
	if i != 1234567 {
		t.Error(intValue, 1234567, i)
	}
	w.Close()
}

func TestUnserializeInt(t *testing.T) {
	intValue := "1234567"
	u := uint(math.MaxUint64)
	data := map[interface{}]int{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		-100:                             -100,
		math.MinInt32:                    int(math.MinInt32),
		math.MaxInt64:                    int(math.MaxInt64),
		math.MinInt64:                    int(math.MinInt64),
		u:                                int(u),
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&intValue:                        1234567,
		time.Unix(123, 456):              123000000456,
		time.Unix(1234567890, 123456789): 1234567890123456789,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&intValue)
	reader := NewReader(w.Bytes(), false)
	var p int
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != 1234567 {
		t.Error(intValue, 1234567, p)
	}
	w.Close()
}

func BenchmarkReadInt(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(12345)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadInt()
	}
	w.Close()
}

func BenchmarkUnserializeInt(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(12345)
	bytes := w.Bytes()
	var p int
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}

func TestReadUint(t *testing.T) {
	intValue := "1234567"
	u := uint(math.MaxUint64)
	data := map[interface{}]uint64{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		math.MaxInt64:                    uint64(math.MaxInt64),
		u:                                uint64(u),
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&intValue:                        1234567,
		time.Unix(123, 456):              123000000456,
		time.Unix(1234567890, 123456789): 1234567890123456789,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&intValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		i := reader.ReadUint()
		if i != data[k] {
			t.Error(k, data[k], i)
		}
	}
	i := reader.ReadUint()
	if i != 1234567 {
		t.Error(intValue, 1234567, i)
	}
	w.Close()
}

func TestUnserializeUint(t *testing.T) {
	intValue := "1234567"
	u := uint(math.MaxUint64)
	data := map[interface{}]uint{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		math.MaxInt64:                    uint(math.MaxInt64),
		u:                                uint(u),
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&intValue:                        1234567,
		time.Unix(123, 456):              123000000456,
		time.Unix(1234567890, 123456789): 1234567890123456789,
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&intValue)
	reader := NewReader(w.Bytes(), false)
	var p uint
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != 1234567 {
		t.Error(intValue, 1234567, p)
	}
	w.Close()
}

func BenchmarkReadUint(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(12345)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadUint()
	}
	w.Close()
}

func BenchmarkUnserializeUint(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(12345)
	bytes := w.Bytes()
	var p uint
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}

func TestReadFloat32(t *testing.T) {
	floatValue := "3.14159"
	data := map[interface{}]float32{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		math.MaxInt64:                    float32(math.MaxInt64),
		math.MaxFloat32:                  math.MaxFloat32,
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&floatValue:                      3.14159,
		time.Unix(123, 456):              float32(123.000000456),
		time.Unix(1234567890, 123456789): float32(1234567890.123456789),
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&floatValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		x := reader.ReadFloat32()
		if x != data[k] {
			t.Error(k, data[k], x)
		}
	}
	x := reader.ReadFloat32()
	if x != float32(3.14159) {
		t.Error(floatValue, 3.14159, x)
	}
	w.Close()
}

func TestUnserializeFloat32(t *testing.T) {
	floatValue := "3.14159"
	data := map[interface{}]float32{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		math.MaxInt64:                    float32(math.MaxInt64),
		math.MaxFloat32:                  math.MaxFloat32,
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&floatValue:                      3.14159,
		time.Unix(123, 456):              float32(123.000000456),
		time.Unix(1234567890, 123456789): float32(1234567890.123456789),
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&floatValue)
	reader := NewReader(w.Bytes(), false)
	var p float32
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != float32(3.14159) {
		t.Error(floatValue, 3.14159, p)
	}
	w.Close()
}

func BenchmarkReadFloat32(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(3.14159)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadFloat32()
	}
	w.Close()
}

func BenchmarkUnserializeFloat32(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(3.14159)
	bytes := w.Bytes()
	var p float32
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}

func TestReadFloat64(t *testing.T) {
	floatValue := "3.14159"
	data := map[interface{}]float64{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		math.MaxFloat32:                  math.MaxFloat32,
		math.MaxFloat64:                  math.MaxFloat64,
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&floatValue:                      3.14159,
		time.Unix(123, 456):              float64(123.000000456),
		time.Unix(1234567890, 123456789): float64(1234567890.123456789),
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&floatValue)
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		x := reader.ReadFloat64()
		if x != data[k] {
			t.Error(k, data[k], x)
		}
	}
	x := reader.ReadFloat64()
	if x != float64(3.14159) {
		t.Error(floatValue, 3.14159, x)
	}
	w.Close()
}

func TestUnserializeFloat64(t *testing.T) {
	floatValue := "3.14159"
	data := map[interface{}]float64{
		true:                             1,
		false:                            0,
		nil:                              0,
		"":                               0,
		0:                                0,
		1:                                1,
		9:                                9,
		100:                              100,
		math.MaxFloat32:                  math.MaxFloat32,
		math.MaxFloat64:                  math.MaxFloat64,
		0.0:                              0,
		"1":                              1,
		"9":                              9,
		&floatValue:                      3.14159,
		time.Unix(123, 456):              float64(123.000000456),
		time.Unix(1234567890, 123456789): float64(1234567890.123456789),
	}
	w := NewWriter(false)
	keys := []interface{}{}
	for k := range data {
		w.Serialize(k)
		keys = append(keys, k)
	}
	w.Serialize(&floatValue)
	reader := NewReader(w.Bytes(), false)
	var p float64
	for _, k := range keys {
		reader.Unserialize(&p)
		if p != data[k] {
			t.Error(k, data[k], p)
		}
	}
	reader.Unserialize(&p)
	if p != float64(3.14159) {
		t.Error(floatValue, 3.14159, p)
	}
	w.Close()
}

func BenchmarkReadFloat64(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(3.14159)
	bytes := w.Bytes()
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.ReadFloat64()
	}
	w.Close()
}

func BenchmarkUnserializeFloat64(b *testing.B) {
	w := NewWriter(true)
	w.Serialize(3.14159)
	bytes := w.Bytes()
	var p float64
	for i := 0; i < b.N; i++ {
		reader := NewReader(bytes, true)
		reader.Unserialize(&p)
	}
	w.Close()
}
