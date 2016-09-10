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
 * io/formatter.go                                        *
 *                                                        *
 * io Formatter for Go.                                   *
 *                                                        *
 * LastModified: Sep 10, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import "testing"

func TestFormat(t *testing.T) {
	var i int
	Unmarshal(Marshal(123), &i)
	if i != 123 {
		t.Error(i)
	}
}

func BenchmarkMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := Marshal(123)
		Recycle(buf)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	buf := Marshal(123)
	var x int
	for i := 0; i < b.N; i++ {
		Unmarshal(buf, &x)
	}
	Recycle(buf)
}
