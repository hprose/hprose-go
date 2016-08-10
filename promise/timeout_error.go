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
 * promise/timeout_error.go                               *
 *                                                        *
 * promise TimeoutError for Go.                           *
 *                                                        *
 * LastModified: Aug 8, 2015                              *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

// TimeoutError is the default error of Promise.Timeout.
type TimeoutError struct{}

// Error implements the TimeoutError Error method.
func (TimeoutError) Error() string {
	return "timeout"
}
