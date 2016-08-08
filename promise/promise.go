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
 * promise for Go.                                        *
 *                                                        *
 * LastModified: Aug 8, 2015                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

// Promise is an implementation of the JS Promise/A+ spec
// (https://promisesaplus.com/).
type Promise struct {
	value  interface{}
	reason error
	state  State
}

// IsPromise used to determine whether the obj is Promise
func IsPromise(obj interface{}) (ok bool) {
	_, ok = obj.(Promise)
	if ok {
		return true
	}
	_, ok = obj.(*Promise)
	return ok
}
