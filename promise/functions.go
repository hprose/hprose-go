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
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"errors"
	"reflect"
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
func Create(computation Computation) Promise {
	promise := New()
	go call(promise, reflect.ValueOf(computation), reflect.Value{})
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
func Sync(computation Computation) Promise {
	promise := New()
	call(promise, reflect.ValueOf(computation), reflect.Value{})
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
		computation := reflect.ValueOf(value)
		if computation.Kind() == reflect.Func {
			call(promise, computation, reflect.Value{})
		} else {
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
		ToPromise(value).Then(
			allHandler(promise, result, index, &count),
			promise.Reject)
	}
	return promise
}

func allHandler(promise Promise, result []interface{}, index int, count *int64) func(value interface{}) {
	return func(value interface{}) {
		result[index] = value
		if atomic.AddInt64(count, -1) == 0 {
			promise.Resolve(result)
		}
	}
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
// Or reject with a IllegalArgumentError if the input iterable is empty--i.e. it
// is impossible to have one winner.
//
// Or reject with an iterable of all the rejection reasons, if the input
// iterable is non-empty, and all input promises reject.
func Any(iterable ...interface{}) Promise {
	count := int64(len(iterable))
	if count == 0 {
		return Reject(IllegalArgumentError("Any(): iterable must not be empty"))
	}
	promise := New()
	for _, value := range iterable {
		ToPromise(value).Then(promise.Resolve, anyHandler(promise, &count))
	}
	return promise
}

func anyHandler(promise Promise, count *int64) func() {
	return func() {
		if atomic.AddInt64(count, -1) == 0 {
			promise.Reject(errors.New("Any(): all promises failed"))
		}
	}
}

// Each function executes a provided function once per element of iterable.
//
// The callback parameter is a function to execute for each element:
//
//		func(value interface{})
//
// value: The current element being processed.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
func Each(callback func(interface{}), iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) {
		if a != nil {
			each(callback, a.([]interface{}))
		}
	})
}

func each(callback func(interface{}), iterable []interface{}) {
	for _, value := range iterable {
		callback(value)
	}
}

// Every function tests whether all elements in the iterable pass the test
// implemented by the provided function.
//
// The callback parameter is a function to test for each element:
//
//		func(value interface{}) bool
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
func Every(callback func(interface{}) bool, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) bool {
		if a == nil {
			return true
		}
		return every(callback, a.([]interface{}))
	})
}

func every(callback func(interface{}) bool, iterable []interface{}) bool {
	for _, value := range iterable {
		if !callback(value) {
			return false
		}
	}
	return true
}

// Some function tests whether some element in the iterable passes the test
// implemented by the provided function.
//
// The callback parameter is a function to test for each element:
//
//		func(value interface{}) bool
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
func Some(callback func(interface{}) bool, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) bool {
		if a == nil {
			return false
		}
		return some(callback, a.([]interface{}))
	})
}

func some(callback func(interface{}) bool, iterable []interface{}) bool {
	for _, value := range iterable {
		if callback(value) {
			return true
		}
	}
	return false
}

// Filter function returns a promise fulfill with a slice that has all elements
// pass the test implemented by the provided function.
//
// The callback parameter is a function to test for each element:
//
//		func(value interface{}) bool
//
// value: The current element being processed.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
func Filter(callback func(interface{}) bool, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) []interface{} {
		if a == nil {
			return nil
		}
		return filter(callback, a.([]interface{}))
	})
}

func filter(callback func(interface{}) bool, iterable []interface{}) []interface{} {
	result := make([]interface{}, 0, len(iterable))
	for _, value := range iterable {
		if callback(value) {
			result = append(result, value)
		}
	}
	return result
}

// Map function returns a promise fulfill with a slice that has the results of
// calling a provided function on every element in iterable.
//
// The callback parameter produces an element of the new slice:
//
//		func(value interface{}) interface{}
//
// value: The current element being processed.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
func Map(callback func(interface{}) interface{}, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) []interface{} {
		if a == nil {
			return nil
		}
		iterable := a.([]interface{})
		result := make([]interface{}, len(iterable))
		for index, value := range iterable {
			result[index] = callback(value)
		}
		return result
	})
}

// Reduce function applies a function against an accumulator and each value of
// the iterable (from left-to-right) to reduce it to a single value.
//
// The callback parameter executes on each value in the iterable, taking three
// arguments:
//
//     func(prev interface{}, value interface{}) interface{}
//
// prev: The value previously returned in the last invocation of the callback.
//
// value: The current element being processed.
//
// The first time the callback is called, prev will be equal to the first value
// in the iterable and value will be equal to the second.
//
// If any of the promises in iterable is rejected, the callback will not be
// executed. the returned promise will be rejected with the rejection reason
// of the first promise that was rejected.
//
// If iterable is empty, the returned promise will be rejected with a
// IllegalArgumentError.
func Reduce(callback func(interface{}, interface{}) interface{}, iterable ...interface{}) Promise {
	return All(iterable...).Then(func(a interface{}) (interface{}, error) {
		if a == nil {
			return nil, IllegalArgumentError("Reduce(): iterable must not be empty")
		}
		return reduce(callback, a.([]interface{})), nil
	})
}

func reduce(callback func(interface{}, interface{}) interface{}, iterable []interface{}) interface{} {
	count := len(iterable)
	result := iterable[0]
	for index := 1; index < count; index++ {
		result = callback(result, iterable[index])
	}
	return result
}
