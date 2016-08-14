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
 * LastModified: Aug 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"reflect"
	"time"
)

type func0 func() (interface{}, error)
type func1 func()
type func2 func(interface{}) (interface{}, error)
type func3 func(interface{})
type func4 func(error) (interface{}, error)
type func5 func(error)

func catch(promise Promise) {
	if e := recover(); e != nil {
		promise.Reject(NewPanicError(e))
	}
}

func call(promise Promise, computation func0) {
	defer catch(promise)
	if result, err := computation(); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func call1(promise Promise, computation func1) {
	defer catch(promise)
	computation()
	promise.Resolve(nil)
}

func call2(promise Promise, computation func2, x interface{}) {
	defer catch(promise)
	if result, err := computation(x); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func call3(promise Promise, computation func3, x interface{}) {
	defer catch(promise)
	computation(x)
	promise.Resolve(nil)
}

func call4(promise Promise, computation func4, e error) {
	defer catch(promise)
	if result, err := computation(e); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func call5(promise Promise, computation func5, e error) {
	defer catch(promise)
	computation(e)
	promise.Resolve(nil)
}

func resolve(next Promise, onFulfilled OnFulfilled, x interface{}) {
	switch f := onFulfilled.(type) {
	case nil:
		next.Resolve(x)
	case func() (interface{}, error):
		go call(next, f)
	case func():
		go call1(next, f)
	case func(interface{}) (interface{}, error):
		go call2(next, f, x)
	case func(interface{}):
		go call3(next, f, x)
	default:
		panic("onFulfilled can't support this type: " +
			reflect.TypeOf(onFulfilled).Name())
	}
}

func reject(next Promise, onRejected OnRejected, e error) {
	switch f := onRejected.(type) {
	case nil:
		next.Reject(e)
	case func() (interface{}, error):
		go call(next, f)
	case func():
		go call1(next, f)
	case func(interface{}) (interface{}, error):
		go call2(next, f, e)
	case func(interface{}):
		go call3(next, f, e)
	case func(error) (interface{}, error):
		go call4(next, f, e)
	case func(error):
		go call5(next, f, e)
	default:
		panic("onRejected can't support this type: " +
			reflect.TypeOf(onRejected).Name())
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
