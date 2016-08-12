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
 * LastModified: Aug 12, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "testing"

func testNew1(t *testing.T) {
	p := New()
	if p.State() != PENDING {
		t.Error("p.State must be PENDING")
	}
}

func testNew2(t *testing.T) {
	p := New(123)
	if p.State() != FULFILLED {
		t.Error("p.State must be FULFILLED")
	}
	if v, _ := p.Get(); v.(int) != 123 {
		t.Error("p.Get value be 123")
	}
}

func TestNew(t *testing.T) {
	testNew1(t)
	testNew2(t)
}
