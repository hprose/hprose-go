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
 * LastModified: Aug 13, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "time"

// Callable is a function type.
// It has no arguments and returns a result with an error.
type Callable func() (interface{}, error)

// OnFulfilled is a function called when the Promise is fulfilled.
// This function has one argument, the fulfillment value.
type OnFulfilled func(interface{}) (interface{}, error)

// OnRejected is a function called when the Promise is rejected.
// This function has one argument, the rejection reason.
type OnRejected func(error) (interface{}, error)

// OnCompleted is a function called when the Promise is completed.
// This function has one argument,
// the fulfillment value when the Promise is fulfilled,
// or the rejection reason when the Promise is rejected.
type OnCompleted func(interface{}) (interface{}, error)

// OnfulfilledSideEffect is a function used as the argument of Promise.Tap
type OnfulfilledSideEffect func(interface{})

// TestFunc is a function used as the argument of Promise.Catch
type TestFunc func(error) bool

// Thenable is an interface that defines a Then method.
type Thenable interface {
	// Then method returns a Promise. It takes two arguments: callback functions
	// for the success and failure cases of the Promise.
	Then(onFulfilled OnFulfilled, onRejected ...OnRejected) Promise
}

// Promise is an interface of the JS Promise/A+ spec
// (https://promisesaplus.com/).
type Promise interface {
	Thenable

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
	Catch(onRejected OnRejected, test ...TestFunc) Promise

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
	// value. If the value is a Thenable (i.e. has a Then method), the returned
	// promise will "follow" that Thenable, adopting its eventual state;
	// otherwise the returned promise will be fulfilled with the value.
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
	Tap(onfulfilledSideEffect OnfulfilledSideEffect) Promise

	// Get the value and reason synchronously, if this promise in PENDING state.
	// this method will block the current goroutine.
	Get() (interface{}, error)
}

func catch(promise Promise) {
	if e := recover(); e != nil {
		promise.Reject(NewPanicError(e))
	}
}

func call(promise Promise, computation Callable) {
	defer catch(promise)
	if result, err := computation(); err != nil {
		promise.Reject(err)
	} else {
		promise.Resolve(result)
	}
}

func resolve(next Promise, onFulfilled OnFulfilled, x interface{}) {
	if onFulfilled != nil {
		go call(next, func() (interface{}, error) { return onFulfilled(x) })
	} else {
		next.Resolve(x)
	}
}

func reject(next Promise, onRejected OnRejected, e error) {
	if onRejected != nil {
		go call(next, func() (interface{}, error) { return onRejected(e) })
	} else {
		next.Reject(e)
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

func tap(promise Promise, onfulfilledSideEffect OnfulfilledSideEffect) Promise {
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
func Create(computation Callable) Promise {
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
func Sync(computation Callable) Promise {
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
		if computation, ok := value.(Callable); ok {
			call(promise, computation)
		} else {
			promise.Resolve(value)
		}
	}()
	return promise
}
