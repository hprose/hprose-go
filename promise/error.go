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
 * LastModified: Aug 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"fmt"
	"runtime"
)

// TimeoutError represents an error when an operation times out.
type TimeoutError struct{}

// Error implements the TimeoutError Error method.
func (TimeoutError) Error() string {
	return "timeout"
}

// TypeError represents an error when a value is not of the expected type.
type TypeError struct {
	message string
}

// Error implements the TypeError Error method.
func (e TypeError) Error() string {
	return e.message
}

// PanicError represents a panic error
type PanicError struct {
	Panic interface{}
	Stack []byte
}

func stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return buf[:n]
		}
		buf = make([]byte, 2*len(buf))
	}
}

// NewPanicError return a panic error
func NewPanicError(v interface{}) *PanicError {
	return &PanicError{v, stack()}
}

// Error implements the PanicError Error method.
func (pe PanicError) Error() string {
	return fmt.Sprintf("%v", pe.Panic)
}
