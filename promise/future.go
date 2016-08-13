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
 * promise/future.go                                      *
 *                                                        *
 * future promise implementation for Go.                  *
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

type future struct {
	value       interface{}
	reason      error
	state       uint32
	subscribers []subscriber
}

// New creates a PENDING Promise object
func New() Promise {
	return new(future)
}

func (p *future) then(onFulfilled OnFulfilled, onRejected OnRejected) Promise {
	next := new(future)
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

func (p *future) Then(onFulfilled OnFulfilled, onRejected ...OnRejected) Promise {
	if len(onRejected) == 0 {
		return p.then(onFulfilled, nil)
	}
	return p.then(onFulfilled, onRejected[0])
}

func (p *future) catch(onRejected OnRejected, test TestFunc) Promise {
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

func (p *future) Catch(onRejected OnRejected, test ...TestFunc) Promise {
	if len(test) == 0 {
		return p.catch(onRejected, nil)
	}
	return p.catch(onRejected, test[0])
}

func (p *future) Complete(onCompleted OnCompleted) Promise {
	return p.then(OnFulfilled(onCompleted), func(e error) (interface{}, error) {
		return onCompleted(e)
	})
}

func (p *future) WhenComplete(action func()) Promise {
	return p.then(func(v interface{}) (interface{}, error) {
		action()
		return v, nil
	}, func(e error) (interface{}, error) {
		action()
		return nil, e
	})
}

func (p *future) Done(onFulfilled OnFulfilled, onRejected ...OnRejected) {
	p.
		Then(onFulfilled, onRejected...).
		Then(nil, func(e error) (interface{}, error) {
			go panic(e)
			return nil, nil
		})
}

func (p *future) Fail(onRejected OnRejected) {
	p.Done(nil, onRejected)
}

func (p *future) Always(onCompleted OnCompleted) {
	p.Done(OnFulfilled(onCompleted), func(e error) (interface{}, error) {
		return onCompleted(e)
	})
}

func (p *future) State() State {
	return State(p.state)
}

func (p *future) resolveThenable(thenable Thenable) {
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

func (p *future) reslove(value interface{}) {
	if atomic.CompareAndSwapUint32(&p.state, uint32(PENDING), uint32(FULFILLED)) {
		p.value = value
		subscribers := p.subscribers
		p.subscribers = nil
		for _, subscriber := range subscribers {
			resolve(subscriber.onFulfilled, subscriber.next, value)
		}
	}
}

func (p *future) Resolve(value interface{}) {
	if promise, ok := value.(*future); ok && promise == p {
		p.Reject(TypeError{"Self resolution"})
	} else if promise, ok := value.(Promise); ok {
		promise.Fill(p)
	} else if thenable, ok := value.(Thenable); ok {
		p.resolveThenable(thenable)
	} else {
		p.reslove(value)
	}
}

func (p *future) Reject(reason error) {
	if atomic.CompareAndSwapUint32(&p.state, uint32(PENDING), uint32(REJECTED)) {
		p.reason = reason
		subscribers := p.subscribers
		p.subscribers = nil
		for _, subscriber := range subscribers {
			reject(subscriber.onRejected, subscriber.next, reason)
		}
	}
}

func (p *future) Fill(promise Promise) {
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

func (p *future) Timeout(duration time.Duration, reason ...error) Promise {
	promise := new(future)
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

func (p *future) Delay(duration time.Duration) Promise {
	promise := new(future)
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

func (p *future) Tap(onfulfilledSideEffect OnfulfilledSideEffect) Promise {
	return p.then(func(v interface{}) (interface{}, error) {
		onfulfilledSideEffect(v)
		return v, nil
	}, nil)
}

func (p *future) Get() (interface{}, error) {
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
