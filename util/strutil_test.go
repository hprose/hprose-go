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
 * util/strutil_test.go                                   *
 *                                                        *
 * strutil test for Go.                                   *
 *                                                        *
 * LastModified: Aug 22, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package util

import (
	"strconv"
	"testing"
)

func TestUTF16Length(t *testing.T) {
	data := map[string]int{
		"":           0,
		"你":          1,
		"你好":         2,
		"你好啊,hello!": 10,
		"🇨🇳":         4,
	}
	for k, v := range data {
		if UTF16Length(k) != v {
			t.Error("The UTF16Length of \"" + k + "\" must be " + strconv.Itoa(v))
		}
	}
}
