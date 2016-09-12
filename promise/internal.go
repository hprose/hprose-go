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
 * promise/internal.go                                    *
 *                                                        *
 * some internal type & functions.                        *
 *                                                        *
 * LastModified: Sep 12, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"reflect"
	"time"
)

func catch(promise Promise) {
	if e := recover(); e != nil {
		promise.Reject(NewPanicError(e))
	}
}

func afterCall(promise Promise, results []reflect.Value) {
	switch len(results) {
	case 0:
		promise.Resolve(nil)
	case 1:
		if results[0].IsNil() {
			promise.Resolve(nil)
		} else {
			result := results[0].Interface()
			if reason := result.(error); reason != nil {
				promise.Reject(reason)
			} else {
				promise.Resolve(result)
			}
		}
	case 2:
		if reason := results[1].Interface().(error); reason != nil {
			promise.Reject(reason)
		} else if results[0].IsNil() {
			promise.Resolve(nil)
		} else {
			promise.Resolve(results[0].Interface())
		}
	}
}

func call(promise Promise, computation reflect.Value, x reflect.Value) {
	defer catch(promise)
	typ := computation.Type()
	numin := typ.NumIn()
	numout := typ.NumOut()
	if numout > 2 {
		panic("The out parameters of computation can't be more than 2")
	}
	switch numin {
	case 0:
		afterCall(promise, computation.Call([]reflect.Value{}))
	case 1:
		afterCall(promise, computation.Call([]reflect.Value{x}))
	default:
		panic("The in parameters of computation can't be more than 1")
	}
}

func resolve(next Promise, onFulfilled OnFulfilled, x interface{}) {
	if onFulfilled == nil {
		next.Resolve(x)
	} else {
		f := reflect.ValueOf(onFulfilled)
		if f.Kind() != reflect.Func {
			panic("onFulfilled can't support this type: " +
				f.Type().String())
		}
		call(next, f, reflect.ValueOf(x))
	}
}

func reject(next Promise, onRejected OnRejected, e error) {
	if onRejected == nil {
		next.Reject(e)
	} else {
		f := reflect.ValueOf(onRejected)
		if f.Kind() != reflect.Func {
			panic("onRejected can't support this type: " +
				f.Type().String())
		}
		call(next, f, reflect.ValueOf(e))
	}
}

func timeout(promise Promise, duration time.Duration, reason ...error) Promise {
	next := New()
	timer := time.AfterFunc(duration, func() {
		if len(reason) > 0 {
			next.Reject(reason[0])
		} else {
			next.Reject(TimeoutError{})
		}
	})
	promise.WhenComplete(func() { timer.Stop() }).Fill(next)
	return next
}

func tap(promise Promise, onfulfilledSideEffect func(interface{})) Promise {
	return promise.Then(func(v interface{}) (interface{}, error) {
		onfulfilledSideEffect(v)
		return v, nil
	})
}
