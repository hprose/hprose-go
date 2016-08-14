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
 * LastModified: Aug 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"errors"
	"sync/atomic"

	"github.com/hprose/hprose-golang/promise"
)

func allHandler(promise Promise, count *int64, result []interface{}, value interface{}, index int) {
	ToPromise(value).Then(func(value interface{}) {
		result[index] = value
		if atomic.AddInt64(count, -1) == 0 {
			promise.Resolve(result)
		}
	}, promise.Reject)
}

// All function returns a promise that resolves when all of the promises in the
// iterable argument have resolved, or rejects with the reason of the first
// passed promise that rejects.
func All(iterable ...interface{}) Promise {
	count := int64(len(iterable))
	result := make([]interface{}, count)
	if count == 0 {
		return Resolve(result)
	}
	promise := New()
	for index, value := range iterable {
		allHandler(promise, &count, result, value, index)
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

func getAnyOnRejected(count *int64) {
	return func() {
		if atomic.AddInt64(count, -1) == 0 {
			promise.Reject(errors.New("any(): all promises failed"))
		}
	}
}

// Any function is a competitive race that allows one winner.
//
// The returned promise will:
//
//     fulfill as soon as any one of the input promises fulfills, with the
//     value of the fulfilled input promise, or
//
//     reject:
//
//         with a IllegalArgumentError if the input array is empty--i.e.
//		   it is impossible to have one winner.
//
//         with an array of all the rejection reasons, if the input array is
//         non-empty, and all input promises reject.
func Any(iterable ...interface{}) Promise {
	count := int64(len(iterable))
	if count == 0 {
		return Reject(IllegalArgumentError("any(): array must not be empty"))
	}
	promise := New()
	for _, value := range iterable {
		ToPromise(value).Then(promise.Resolve, getAnyOnRejected(&count))
	}
	return promise
}
