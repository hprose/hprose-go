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
 * promise/promise_func.go                                *
 *                                                        *
 * promise functions for Go.                              *
 *                                                        *
 * LastModified: Aug 13, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "time"

// New creates a Promise object
func New(value ...interface{}) Promise {
	promise := new(promiseImpl)
	if len(value) > 0 {
		if computation, ok := value[0].(Callable); ok {
			go promise.init(computation)
		} else if e, ok := value[0].(error); ok {
			promise.Reject(e)
		} else {
			promise.Resolve(value[0])
		}
	}
	return promise
}

// Sync creates a Promise object containing the result of immediately calling
// computation.
//
// If calling computation returns error, the returned Promise is rejected with
// the error.
//
// If calling computation returns a Promise object, completion of the created
// Promise will wait until the returned Promise completes, and will then
// complete with the same result.
//
// If calling computation returns a non-Promise value, the returned Promise is
// completed with that value.
func Sync(computation Callable) Promise {
	promise := new(promiseImpl)
	promise.init(computation)
	return promise
}

// Delayed creates a Promise object with the given value after a delay.
//
// If the value is a Callable function, it will be executed after the given duration has passed, and the Promise object is completed with the result.
func Delayed(duration time.Duration, value interface{}) Promise {
	promise := new(promiseImpl)
	go func() {
		time.Sleep(duration)
		if computation, ok := value.(Callable); ok {
			promise.init(computation)
		} else if e, ok := value.(error); ok {
			promise.Reject(e)
		} else {
			promise.Resolve(value)
		}
	}()
	return promise
}
