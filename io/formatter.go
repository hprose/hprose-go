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

// Serialize data
func Serialize(v interface{}, simple bool) []byte {
	return NewWriter(simple).Serialize(v).Bytes()
}

// Marshal data
func Marshal(v interface{}) []byte {
	return Serialize(v, true)
}

// Unserialize data
func Unserialize(b []byte, p interface{}, simple bool) {
	NewReader(b, simple).Unserialize(p)
}

// Unmarshal data
func Unmarshal(b []byte, p interface{}) {
	Unserialize(b, p, true)
}
