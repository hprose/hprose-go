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
 * promise/promise.go                                     *
 *                                                        *
 * promise interface for Go.                              *
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

// OnFulfilled is a function called when the Promise is fulfilled.
// This function has one argument, the fulfillment value.
//
// The function type can be the following:
//     func() (interface{}, error)
//     func()
//     func(interface{}) (interface{}, error)
//     func(interface{})
type OnFulfilled interface{}

// OnRejected is a function called when the Promise is rejected.
// This function has one argument, the rejection reason.
//
// The function type can be the following:
//     func() (interface{}, error)
//     func()
//     func(interface{}) (interface{}, error)
//     func(interface{})
//     func(error) (interface{}, error)
//     func(error)
type OnRejected interface{}

// OnCompleted is a function called when the Promise is completed.
// This function has one argument,
// the fulfillment value when the Promise is fulfilled,
// or the rejection reason when the Promise is rejected.
// The function type can be the following:
//     func() (interface{}, error)
//     func()
//     func(interface{}) (interface{}, error)
//     func(interface{})
type OnCompleted interface{}

// Promise is an interface of the JS Promise/A+ spec
// (https://promisesaplus.com/).
type Promise interface {
	// Then method returns a Promise. It takes two arguments: callback functions
	// for the success and failure cases of the Promise.
	Then(onFulfilled OnFulfilled, onRejected ...OnRejected) Promise

	// Catch handles errors emitted by this Promise.
	//
	// This is the asynchronous equivalent of a "catch" block.
	//
	// Returns a new Promise that will be completed with either the result of
	// this promise or the result of calling the onRejected callback.
	//
	// If this promise completes with a value, the returned promise completes
	// with the same value.
	//
	// If this promise completes with an error, then test is first called with
	// the error value.
	//
	// If test returns false, the error is not handled by this Catch, and the
	// returned promise completes with the same error and stack trace as this
	// promise.
	//
	// If test returns true, onRejected is called with the error and possibly
	// stack trace, and the returned promise is completed with the result of
	// this call in exactly the same way as for Then's onRejected.
	//
	// If test is omitted, it defaults to a function that always returns true.
	// The test function should not panic, but if it does, it is handled as if
	// the the onRejected function had panic.
	Catch(onRejected OnRejected, test ...func(error) bool) Promise

	// Complete is the same way as Then(onCompleted, onCompleted)
	Complete(onCompleted OnCompleted) Promise

	// WhenComplete register a function to be called when the promise completes.
	//
	// The action function is called when this promise completes, whether it
	// does so with a value or with an error.
	//
	// If this promise completes with a value, the returned promise completes
	// with the same value.
	//
	// If this promise completes with an error, the returned promise completes
	// with the same error.
	//
	// The action function should not panic, but if it does, the returned
	// promise completes with a PanicError.
	WhenComplete(action func()) Promise

	// Done is the same semantics as Then except that it don't return a Promise.
	// If the callback function (onFulfilled or onRejected) returns error or
	// panics, the application will be crashing.
	// The result of the callback function will be ignored.
	Done(onFulfilled OnFulfilled, onRejected ...OnRejected)

	// State return the current state of the Promise
	State() State

	// Resolve method returns a Promise object that is resolved with the given
	// value. If the value is a Promise, the returned promise will "follow" that Promise, adopting its eventual state; otherwise the returned promise
	// will be fulfilled with the value.
	Resolve(value interface{})

	// Reject method returns a Promise object that is rejected with the given
	// reason.
	Reject(reason error)

	// Fill the promise with this promise if the promise is in PENDING state.
	// otherwise nothing to do.
	Fill(promise Promise)

	// Timeout create a new promise that will reject with a TimeoutError or a
	// custom reason after a timeout if promise does not fulfill or reject
	// beforehand.
	Timeout(duration time.Duration, reason ...error) Promise

	// Delay create a new promise that will, after duration delay, fulfill with
	// the same value as this promise. If this promise rejects, delayed promise
	// will be rejected immediately.
	Delay(duration time.Duration) Promise

	// Tap executes a function as a side effect when promise fulfills.
	//
	// It returns a new promise:
	// 1. If promise fulfills, onFulfilledSideEffect is executed:
	//     * If onFulfilledSideEffect returns successfully, the promise
	//       returned by tap fulfills with promise's original fulfillment
	//       value.
	//     * If onFulfilledSideEffect panics, the promise returned by tap
	//       rejects with the panic message as the reason.
	// 2. If promise rejects, onFulfilledSideEffect is not executed, and the
	//    promise returned by tap rejects with promise's rejection reason.
	Tap(onfulfilledSideEffect func(interface{})) Promise

	// Get the value and reason synchronously, if this promise in PENDING state.
	// this method will block the current goroutine.
	Get() (interface{}, error)
}

func catch(promise Promise) {
	if e := recover(); e != nil {
		promise.Reject(NewPanicError(e))
	}
}

func call(promise Promise, computation func() (interface{}, error)) {
	defer catch(promise)
	if result, err := computation(); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func call1(promise Promise, computation func()) {
	defer catch(promise)
	computation()
	promise.Resolve(nil)
}

func call2(promise Promise, computation func(interface{}) (interface{}, error), x interface{}) {
	defer catch(promise)
	if result, err := computation(x); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func call3(promise Promise, computation func(interface{}), x interface{}) {
	defer catch(promise)
	computation(x)
	promise.Resolve(nil)
}

func call4(promise Promise, computation func(error) (interface{}, error), e error) {
	defer catch(promise)
	if result, err := computation(e); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func call5(promise Promise, computation func(error), e error) {
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
		panic("onFulfilled can't support this type: " + reflect.TypeOf(onFulfilled).Name())
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
		panic("onRejected can't support this type: " + reflect.TypeOf(onRejected).Name())
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
