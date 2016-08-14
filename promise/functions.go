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
)

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
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
//
// Parameters:
//   callback: Function to execute for each element, taking two arguments:
//     index: The index of the current element being processed.
//     value: The current element being processed.
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
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
//
// Parameters:
//   callback: Function to test for each element, taking two arguments:
//     index: The index of the current element being processed.
//     value: The current element being processed.
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
