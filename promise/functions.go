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

import "sync/atomic"

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
