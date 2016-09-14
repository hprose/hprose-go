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
 * rpc/error.go                                           *
 *                                                        *
 * rpc error for Go.                                      *
 *                                                        *
 * LastModified: Sep 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

import (
	"errors"
	"fmt"
	"runtime"
)

var errServerIsAlreadyStarted = errors.New("The server is already started")
var errServerIsNotStarted = errors.New("The server is not started")

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
func (pe *PanicError) Error() string {
	return fmt.Sprintf("%v", pe.Panic)
}
