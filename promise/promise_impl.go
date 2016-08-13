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
 * LastModified: Aug 13, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import (
	"sync/atomic"
	"time"
)

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

// New creates a PENDING Promise object
func New() Promise {
	return new(promiseImpl)
}

func call(onFulfilled OnFulfilled, next Promise, x interface{}) {
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
}

func resolve(onFulfilled OnFulfilled, next Promise, x interface{}) {
	if onFulfilled != nil {
		go call(onFulfilled, next, x)
	} else {
		next.Resolve(x)
	}
}

func reject(onRejected OnRejected, next Promise, e error) {
	if onRejected != nil {
		go call(func(x interface{}) (interface{}, error) {
			return onRejected(x.(error))
		}, next, e)
	} else {
		next.Reject(e)
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
	if test == nil {
		return p.then(nil, onRejected)
	}
	return p.then(nil, func(e error) (interface{}, error) {
		if test(e) {
			return p.then(nil, onRejected), nil
		}
		return nil, e
	})
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

func (p *promiseImpl) resolveThenable(thenable Thenable) {
	var done uint32
	defer func() {
		if e := recover(); e != nil && atomic.CompareAndSwapUint32(&done, 0, 1) {
			p.Reject(NewPanicError(e))
		}
	}()
	thenable.Then(func(y interface{}) (interface{}, error) {
		if atomic.CompareAndSwapUint32(&done, 0, 1) {
			p.Resolve(y)
		}
		return nil, nil
	}, func(e error) (interface{}, error) {
		if atomic.CompareAndSwapUint32(&done, 0, 1) {
			p.Reject(e)
		}
		return nil, nil
	})
}

func (p *promiseImpl) reslove(value interface{}) {
	if atomic.CompareAndSwapUint32(&p.state, uint32(PENDING), uint32(FULFILLED)) {
		p.value = value
		subscribers := p.subscribers
		p.subscribers = nil
		for _, subscriber := range subscribers {
			resolve(subscriber.onFulfilled, subscriber.next, value)
		}
	}
}

func (p *promiseImpl) Resolve(value interface{}) {
	if promise, ok := value.(*promiseImpl); ok && promise == p {
		p.Reject(TypeError{"Self resolution"})
	} else if promise, ok := value.(Promise); ok {
		promise.Fill(p)
	} else if thenable, ok := value.(Thenable); ok {
		p.resolveThenable(thenable)
	} else {
		p.reslove(value)
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
