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
 * util/strutil.go                                        *
 *                                                        *
 * string util for Go.                                    *
 *                                                        *
 * LastModified: Aug 22, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package util

// UTF16Length return the UTF16 length of str.
// str must be an UTF8 encode string, otherwise return -1.
func UTF16Length(str string) (n int) {
	length := len(str)
	n = length
	p := 0
	for p < length {
		a := str[p]
		switch {
		case a < 0x80:
			p++
		case (a & 0xE0) == 0xC0:
			p += 2
			n--
		case (a & 0xF0) == 0xE0:
			p += 3
			n -= 2
		case (a & 0xF8) == 0xF0:
			p += 4
			n -= 2
		default:
			return -1
		}
	}
	return n
}
