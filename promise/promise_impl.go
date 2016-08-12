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
 * promise/promise_impl.go                                *
 *                                                        *
 * promise implementation for Go.                         *
 *                                                        *
 * LastModified: Aug 11, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/
package promise

import (
	"sync/atomic"
	"time"
)

// Callable is a function type.
// It has no arguments and returns a result with an error.
type Callable func() (interface{}, error)

type subscriber struct {
	onFulfilled OnFulfilled
	onRejected  OnRejected
	next        Promise
}

type promiseImpl struct {
	value       interface{}
	reason      error
	state       uint32
	subscribers []subscriber
}

func resolve(onFulfilled OnFulfilled, next Promise, x interface{}) {
	if onFulfilled != nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					next.Reject(NewPanicError(e))
				}
			}()
			if result, err := onFulfilled(x); err != nil {
				next.Reject(err)
			} else {
				next.Resolve(result)
			}
		}()
	} else {
		next.Resolve(x)
	}
}

func reject(onRejected OnRejected, next Promise, e error) {
	if onRejected != nil {
		go func() {
			defer func() {
				if e := recover(); e != nil {
					next.Reject(NewPanicError(e))
				}
			}()
			if result, err := onRejected(e); err != nil {
				next.Reject(err)
			} else {
				next.Resolve(result)
			}
		}()
	} else {
		next.Reject(e)
	}
}

func (p *promiseImpl) init(computation Callable) {
	result, err := computation()
	if err != nil {
		p.Reject(err)
	} else {
		p.Resolve(result)
	}
}

func (p *promiseImpl) then(onFulfilled OnFulfilled, onRejected OnRejected) Promise {
	next := new(promiseImpl)
	switch State(p.state) {
	case FULFILLED:
		resolve(onFulfilled, next, p.value)
	case REJECTED:
		reject(onRejected, next, p.reason)
	default:
		p.subscribers = append(p.subscribers,
			subscriber{onFulfilled, onRejected, next})
	}
	return next
}

func (p *promiseImpl) Then(onFulfilled OnFulfilled, onRejected ...OnRejected) Promise {
	if len(onRejected) == 0 {
		return p.then(onFulfilled, nil)
	}
	return p.then(onFulfilled, onRejected[0])
}

func (p *promiseImpl) catch(onRejected OnRejected, test TestFunc) Promise {
	if test != nil {
		return p.then(nil, func(e error) (interface{}, error) {
			if test(e) {
				return p.then(nil, onRejected), nil
			}
			return nil, e
		})
	}
	return p.then(nil, onRejected)
}

func (p *promiseImpl) Catch(onRejected OnRejected, test ...TestFunc) Promise {
	if len(test) == 0 {
		return p.catch(onRejected, nil)
	}
	return p.catch(onRejected, test[0])
}

func (p *promiseImpl) Complete(onCompleted OnCompleted) Promise {
	return p.then(OnFulfilled(onCompleted), func(e error) (interface{}, error) {
		return onCompleted(e)
	})
}

func (p *promiseImpl) WhenComplete(action func()) Promise {
	return p.then(func(v interface{}) (interface{}, error) {
		action()
		return v, nil
	}, func(e error) (interface{}, error) {
		action()
		return nil, e
	})
}

func (p *promiseImpl) Done(onFulfilled OnFulfilled, onRejected ...OnRejected) {
	p.
		Then(onFulfilled, onRejected...).
		Then(nil, func(e error) (interface{}, error) {
			go panic(e)
			return nil, nil
		})
}

func (p *promiseImpl) Fail(onRejected OnRejected) {
	p.Done(nil, onRejected)
}

func (p *promiseImpl) Always(onCompleted OnCompleted) {
	p.Done(OnFulfilled(onCompleted), func(e error) (interface{}, error) {
		return onCompleted(e)
	})
}

func (p *promiseImpl) State() State {
	return State(p.state)
}

func (p *promiseImpl) resolve(thenable Thenable) {
	var done uint32
	resolveFunc := func(y interface{}) (interface{}, error) {
		if atomic.CompareAndSwapUint32(&done, 0, 1) {
			p.Resolve(y)
		}
		return nil, nil
	}
	rejectFunc := func(e error) (interface{}, error) {
		if atomic.CompareAndSwapUint32(&done, 0, 1) {
			p.Reject(e)
		}
		return nil, nil
	}
	defer func() {
		if e := recover(); e != nil {
			if atomic.CompareAndSwapUint32(&done, 0, 1) {
				p.Reject(NewPanicError(e))
			}
		}
	}()
	thenable.Then(resolveFunc, rejectFunc)
}

func (p *promiseImpl) Resolve(value interface{}) {
	if promise, ok := value.(*promiseImpl); ok && promise == p {
		p.Reject(TypeError{"Self resolution"})
	} else if promise, ok := value.(Promise); ok {
		promise.Fill(p)
	} else if thenable, ok := value.(Thenable); ok {
		p.resolve(thenable)
	} else if atomic.CompareAndSwapUint32(&p.state, uint32(PENDING), uint32(FULFILLED)) {
		p.value = value
		subscribers := p.subscribers
		p.subscribers = nil
		for _, subscriber := range subscribers {
			resolve(subscriber.onFulfilled, subscriber.next, value)
		}
	}
}

func (p *promiseImpl) Reject(reason error) {
	if atomic.CompareAndSwapUint32(&p.state, uint32(PENDING), uint32(REJECTED)) {
		p.reason = reason
		subscribers := p.subscribers
		p.subscribers = nil
		for _, subscriber := range subscribers {
			reject(subscriber.onRejected, subscriber.next, reason)
		}
	}
}

func (p *promiseImpl) Fill(promise Promise) {
	resolveFunc := func(v interface{}) (interface{}, error) {
		promise.Resolve(v)
		return nil, nil
	}
	rejectFunc := func(e error) (interface{}, error) {
		promise.Reject(e)
		return nil, nil
	}
	p.Then(resolveFunc, rejectFunc)
}

func (p *promiseImpl) Timeout(duration time.Duration, reason ...error) Promise {
	promise := new(promiseImpl)
	timer := time.AfterFunc(duration, func() {
		if len(reason) > 0 {
			promise.Reject(reason[0])
		} else {
			promise.Reject(TimeoutError{})
		}
	})
	p.WhenComplete(func() { timer.Stop() }).Fill(promise)
	return promise
}

func (p *promiseImpl) Delay(duration time.Duration) Promise {
	promise := new(promiseImpl)
	p.then(func(v interface{}) (interface{}, error) {
		go func() {
			time.Sleep(duration)
			promise.Resolve(v)
		}()
		return nil, nil
	}, func(e error) (interface{}, error) {
		promise.Reject(e)
		return nil, nil
	})
	return promise
}

func (p *promiseImpl) Tap(onfulfilledSideEffect OnfulfilledSideEffect) Promise {
	return p.then(func(v interface{}) (interface{}, error) {
		onfulfilledSideEffect(v)
		return v, nil
	}, nil)
}

func (p *promiseImpl) Get() (interface{}, error) {
	c := make(chan interface{})
	p.then(func(v interface{}) (interface{}, error) {
		c <- v
		return nil, nil
	}, func(e error) (interface{}, error) {
		c <- e
		return nil, nil
	})
	v := <-c
	if e, ok := v.(error); ok {
		return nil, e
	}
	return v, nil
}

// New creates a Promise object
func New(value ...interface{}) Promise {
	promise := new(promiseImpl)
	if len(value) > 0 {
		if computation, ok := value[0].(Callable); ok {
			go promise.init(computation)
		} else if e, ok := value[0].(error); ok {
			promise.Reject(e)
		} else {
			promise.Resolve(value[0])
		}
	}
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
	promise := new(promiseImpl)
	promise.init(computation)
	return promise
}

// Delayed creates a Promise object with the given value after a delay.
//
// If the value is a Callable function, it will be executed after the given duration has passed, and the Promise object is completed with the result.
func Delayed(duration time.Duration, value interface{}) Promise {
	promise := new(promiseImpl)
	go func() {
		time.Sleep(duration)
		if computation, ok := value.(Callable); ok {
			promise.init(computation)
		} else if e, ok := value.(error); ok {
			promise.Reject(e)
		} else {
			promise.Resolve(value)
		}
	}()
	return promise
}
