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
 * promise/promise_test.go                                *
 *                                                        *
 * promise test for Go.                                   *
 *                                                        *
 * LastModified: Aug 8, 2015                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "testing"

func TestIsPromise(t *testing.T) {
	p := Promise{}
	if !IsPromise(p) {
		t.Error("How is that possible?")
	}
	pp := &Promise{}
	if !IsPromise(pp) {
		t.Error("How is that possible?")
	}
	ppp := &pp
	if IsPromise(ppp) {
		t.Error("How is that possible?")
	}
}
