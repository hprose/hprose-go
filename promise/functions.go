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
 * promise/functions.go                                   *
 *                                                        *
 * some functions of promise for Go.                      *
 *                                                        *
 * LastModified: Aug 13, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package promise

// All function returns a promise that resolves when all of the promises in the
// iterable argument have resolved, or rejects with the reason of the first
// passed promise that rejects.
// func All(iterable ...interface{}) Promise {
// 	if len(iterable) == 0 {
// 		return Resolve(nil)
// 	}
// }
