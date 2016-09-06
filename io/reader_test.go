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

import "testing"

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
	reader := NewReader(w.Bytes(), false)
	for _, k := range keys {
		b := reader.ReadBool()
		if b != data[k] {
			t.Error(k, data[k], b)
		}
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
