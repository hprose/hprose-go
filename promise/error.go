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
 * promise/error.go                                       *
 *                                                        *
 * promise error for Go.                                  *
 *                                                        *
 * LastModified: Aug 18, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"fmt"
	"runtime"

	"github.com/hprose/hprose-golang/pool"
)

// IllegalArgumentError represents an error when a function/method has been
// passed an illegal or inappropriate argument.
type IllegalArgumentError string

// Error implements the IllegalArgumentError Error method.
func (e IllegalArgumentError) Error() string {
	return string(e)
}

// TimeoutError represents an error when an operation times out.
type TimeoutError struct{}

// Error implements the TimeoutError Error method.
func (TimeoutError) Error() string {
	return "timeout"
}

// TypeError represents an error when a value is not of the expected type.
type TypeError string

// Error implements the TypeError Error method.
func (e TypeError) Error() string {
	return string(e)
}

// PanicError represents a panic error
type PanicError struct {
	Panic interface{}
	Stack []byte
}

func stack() []byte {
	size := 1024
	buf := pool.Alloc(size)
	for {
		n := runtime.Stack(buf, false)
		if n < size {
			return buf[:n]
		}
		pool.Recycle(buf)
		size <<= 1
		buf = pool.Alloc(size)
	}
}

// NewPanicError return a panic error
func NewPanicError(v interface{}) *PanicError {
	return &PanicError{v, stack()}
}

// Error implements the PanicError Error method.
func (pe *PanicError) Error() string {
	return fmt.Sprintf("%v", pe.Panic)
}
