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
 * LastModified: Aug 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	p := New()
	if p.State() != PENDING {
		t.Error("p.State must be PENDING")
	}
}

func TestResolve(t *testing.T) {
	p := Resolve(123)
	if p.State() != FULFILLED {
		t.Error("p.State must be FULFILLED")
	}
	if v, _ := p.Get(); v.(int) != 123 {
		t.Error("p.Get value be 123")
	}
}

func TestAll(t *testing.T) {
	p := All(Resolve(1), 2, Resolve(3))
	if v, err := p.Get(); err == nil {
		if fmt.Sprintf("%v", v) != "[1 2 3]" {
			t.Error("v must be [1 2 3]")
		}
	} else {
		t.Error(err)
	}
}
