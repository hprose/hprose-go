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
 * promise/rejected.go                                    *
 *                                                        *
 * rejected promise implementation for Go.                *
 *                                                        *
 * LastModified: Aug 18, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "time"

type rejected struct {
	reason error
}

// Reject creates a REJECTED Promise object
func Reject(reason error) Promise {
	return &rejected{reason}
}

func (p *rejected) Then(onFulfilled OnFulfilled, onRejected ...OnRejected) Promise {
	if len(onRejected) == 0 {
		return &rejected{p.reason}
	}
	next := New()
	reject(next, onRejected[0], p.reason)
	return next
}

func (p *rejected) Catch(onRejected OnRejected, test ...func(error) bool) Promise {
	if len(test) == 0 || test[0](p.reason) {
		next := New()
		reject(next, onRejected, p.reason)
		return next
	}
	return &rejected{p.reason}
}

func (p *rejected) Complete(onCompleted OnCompleted) Promise {
	return p.Then(nil, onCompleted)
}

func (p *rejected) WhenComplete(action func()) Promise {
	return p.Catch(func(e error) (interface{}, error) {
		action()
		return nil, e
	})
}

func (p *rejected) Done(onFulfilled OnFulfilled, onRejected ...OnRejected) {
	p.Then(nil, onRejected...).Then(nil, func(e error) { go panic(e) })
}

func (p *rejected) State() State {
	return REJECTED
}

func (p *rejected) Resolve(value interface{}) {}

func (p *rejected) Reject(reason error) {}

func (p *rejected) Fill(promise Promise) {
	go promise.Reject(p.reason)
}

func (p *rejected) Timeout(duration time.Duration, reason ...error) Promise {
	return timeout(p, duration, reason...)
}

func (p *rejected) Delay(duration time.Duration) Promise {
	return &rejected{p.reason}
}

func (p *rejected) Tap(onfulfilledSideEffect func(interface{})) Promise {
	return &rejected{p.reason}
}

func (p *rejected) Get() (interface{}, error) {
	return nil, p.reason
}
