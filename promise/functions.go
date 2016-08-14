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
 * promise/functions.go                                   *
 *                                                        *
 * some functions of promise for Go.                      *
 *                                                        *
 * LastModified: Aug 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"errors"
	"sync/atomic"
	"time"
)

// Create creates a Promise object containing the result of asynchronously
// calling computation.
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
func Create(computation func() (interface{}, error)) Promise {
	promise := New()
	go call(promise, computation)
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
func Sync(computation func() (interface{}, error)) Promise {
	promise := New()
	call(promise, computation)
	return promise
}

// Delayed creates a Promise object with the given value after a delay.
//
// If the value is a Callable function, it will be executed after the given
// duration has passed, and the Promise object is completed with the result.
func Delayed(duration time.Duration, value interface{}) Promise {
	promise := New()
	go func() {
		time.Sleep(duration)
		switch computation := value.(type) {
		case func() (interface{}, error):
			call(promise, computation)
		case func():
			call1(promise, computation)
		default:
			promise.Resolve(value)
		}
	}()
	return promise
}

// All function returns a promise that resolves when all of the promises in the
// iterable argument have resolved, or rejects with the reason of the first
// passed promise that rejects.
func All(iterable ...interface{}) Promise {
	count := int64(len(iterable))
	if count == 0 {
		return Resolve(nil)
	}
	result := make([]interface{}, count)
	promise := New()
	for index, value := range iterable {
		ToPromise(value).Then(func(index int) func(value interface{}) {
			return func(value interface{}) {
				result[index] = value
				if atomic.AddInt64(&count, -1) == 0 {
					promise.Resolve(result)
				}
			}
		}(index), promise.Reject)
	}
	return promise
}

// Race function returns a promise that resolves or rejects as soon as one of
// the promises in the iterable resolves or rejects, with the value or reason
// from that promise.
func Race(iterable ...interface{}) Promise {
	promise := New()
	for _, value := range iterable {
		ToPromise(value).Fill(promise)
	}
	return promise
}

// Any function is a competitive race that allows one winner.
//
// The returned promise will fulfill as soon as any one of the input promises
// fulfills, with the value of the fulfilled input promise.
//
// Or reject with a IllegalArgumentError if the input array is empty--i.e. it
// is impossible to have one winner.
//
// Or reject with an array of all the rejection reasons, if the input array is
// non-empty, and all input promises reject.
func Any(iterable ...interface{}) Promise {
	count := int64(len(iterable))
	if count == 0 {
		return Reject(IllegalArgumentError("any(): array must not be empty"))
	}
	promise := New()
	for _, value := range iterable {
		ToPromise(value).Then(promise.Resolve, func() {
			if atomic.AddInt64(&count, -1) == 0 {
				promise.Reject(errors.New("any(): all promises failed"))
			}
		})
	}
	return promise
}

// Each function executes a provided function once per element of iterable.
//
// The callback parameter is a function to execute for each element:
//
//		func(index int, value interface{})
//
// index: The index of the current element being processed.
//
// value: The current element being processed.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
func Each(callback func(int, interface{}), iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) {
		if a == nil {
			return
		}
		iterable := a.([]interface{})
		for index, value := range iterable {
			callback(index, value)
		}
	})
}

// Every function tests whether all elements in the iterable pass the test
// implemented by the provided function.
//
// The callback parameter is a function to test for each element:
//
//		func(index int, value interface{}) bool
//
// index: The index of the current element being processed.
//
// value: The current element being processed.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
//
// The returned promise will fulfill with false if the callback function returns
// false for any iterable element immediately. Otherwise, if callback returned
// a true value for all elements, Every will return promise fulfill with true.
//
// If iterable is empty, The returned promise will fulfill with true.
func Every(callback func(int, interface{}) bool, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) (interface{}, error) {
		if a == nil {
			return true, nil
		}
		iterable := a.([]interface{})
		for index, value := range iterable {
			if !callback(index, value) {
				return false, nil
			}
		}
		return true, nil
	})
}

// Some function tests whether some element in the array passes the test
// implemented by the provided function.
//
// The callback parameter is a function to test for each element:
//
//		func(index int, value interface{}) bool
//
// index: The index of the current element being processed.
//
// value: The current element being processed.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
//
// The returned promise will fulfill with true if the callback function returns
// true for any iterable element immediately. Otherwise, if callback returned a
// false value for all elements, Some will return promise fulfill with false.
//
// If iterable is empty, The returned promise will fulfill with false.
func Some(callback func(int, interface{}) bool, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) (interface{}, error) {
		if a == nil {
			return false, nil
		}
		iterable := a.([]interface{})
		for index, value := range iterable {
			if callback(index, value) {
				return true, nil
			}
		}
		return false, nil
	})
}
