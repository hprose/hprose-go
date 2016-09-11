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
 * LastModified: Sep 11, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "time"

// OnFulfilled is a function called when the Promise is fulfilled.
//
// This function has one argument, the fulfillment value.
//
// The function type can be the following form:
//
//     func()
//     func() R
//     func() (R, error)
//     func(T)
//     func(T) R
//     func(T) (R, error)
//
// T and R can be any type, and they are can be the same or different.
type OnFulfilled interface{}

// OnRejected is a function called when the Promise is rejected.
//
// This function has one argument, the rejection reason.
//
// The function type can be the following form:
//
//     func()
//     func() T
//     func() (T, error)
//     func(E)
//     func(E) T
//     func(E) (T, error)
//
// E is error type or interface{}.
// T can be any type.
type OnRejected interface{}

// OnCompleted is a function called when the Promise is completed.
//
// This function has one argument, the fulfillment value when the Promise is
// fulfilled, or the rejection reason when the Promise is rejected.
//
// The function type can be the following form:
//
//     func()
//     func() T
//     func() (T, error)
//     func(interface{})
//     func(interface{}) T
//     func(interface{}) (T, error)
//
// T can be any type.
type OnCompleted interface{}

// Computation is a function for Promise Create.
//
// This function has no argument.
//
// The function type can be the following form:
//
//     func()
//     func() T
//     func() (T, error)
//
// T can be any type.
type Computation interface{}

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
	// value. If the value is a Promise, the returned promise will "follow"
	// that Promise, adopting its eventual state; otherwise the returned promise
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
