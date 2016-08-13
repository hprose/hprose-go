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
 * promise/fulfilled.go                                   *
 *                                                        *
 * fulfilled promise implementation for Go.               *
 *                                                        *
 * LastModified: Aug 13, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

import "time"

type fulfilled struct {
	value interface{}
}

// Resolve creates a Promise object completed with the value
func Resolve(value interface{}) Promise {
	if _, ok := value.(Thenable); ok {
		p := New()
		p.Resolve(value)
		return p
	}
	return fulfilled{value}
}

func (p fulfilled) Then(onFulfilled OnFulfilled, onRejected ...OnRejected) Promise {
	if onFulfilled == nil {
		return fulfilled{p.value}
	}
	next := New()
	resolve(next, onFulfilled, p.value)
	return next
}

func (p fulfilled) Catch(onRejected OnRejected, test ...TestFunc) Promise {
	return fulfilled{p.value}
}

func (p fulfilled) Complete(onCompleted OnCompleted) Promise {
	return p.Then(OnFulfilled(onCompleted))
}

func (p fulfilled) WhenComplete(action func()) Promise {
	return p.Then(func(v interface{}) (interface{}, error) {
		action()
		return v, nil
	})
}

func (p fulfilled) Done(onFulfilled OnFulfilled, onRejected ...OnRejected) {
	p.
		Then(onFulfilled).
		Then(nil, func(e error) (interface{}, error) {
			go panic(e)
			return nil, nil
		})
}

func (p fulfilled) State() State {
	return FULFILLED
}

func (p fulfilled) Resolve(value interface{}) {}

func (p fulfilled) Reject(reason error) {}

func (p fulfilled) Fill(promise Promise) {
	go promise.Resolve(p.value)
}

func (p fulfilled) Timeout(duration time.Duration, reason ...error) Promise {
	return timeout(p, duration, reason...)
}

func (p fulfilled) Delay(duration time.Duration) Promise {
	next := New()
	go func() {
		time.Sleep(duration)
		next.Resolve(p.value)
	}()
	return next
}

func (p fulfilled) Tap(onfulfilledSideEffect OnfulfilledSideEffect) Promise {
	return tap(p, onfulfilledSideEffect)
}

func (p fulfilled) Get() (interface{}, error) {
	return p.value, nil
}
