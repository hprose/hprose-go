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
 * rpc/filter.go                                          *
 *                                                        *
 * hprose filter interface for Go.                        *
 *                                                        *
 * LastModified: Sep 14, 2016                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package rpc

// Filter is hprose filter
type Filter interface {
	InputFilter(data []byte, context Context) []byte
	OutputFilter(data []byte, context Context) []byte
}

// filterManager is the filter manager
type filterManager struct {
	filters []Filter
}

// Filter return the first filter
func (fm *filterManager) Filter() Filter {
	if len(fm.filters) == 0 {
		return nil
	}
	return fm.filters[0]
}

// FilterByIndex return the filter by index
func (fm *filterManager) FilterByIndex(index int) Filter {
	n := len(fm.filters)
	if index < 0 && index >= n {
		return nil
	}
	return fm.filters[index]
}

// SetFilter will replace the current filter settings
func (fm *filterManager) SetFilter(filter ...Filter) {
	fm.filters = make([]Filter, len(filter))
	fm.AddFilter(filter...)
}

// AddFilter add the filter to this FilterManager
func (fm *filterManager) AddFilter(filter ...Filter) {
	if len(filter) > 0 {
		fm.filters = append(fm.filters, filter...)
	}
}

// RemoveFilterByIndex remove the filter by the index
func (fm *filterManager) RemoveFilterByIndex(index int) {
	n := len(fm.filters)
	if index < 0 && index >= n {
		return
	}
	if index == n-1 {
		fm.filters = fm.filters[:index]
	} else {
		fm.filters = append(fm.filters[:index], fm.filters[index+1:]...)
	}
}

func (fm *filterManager) removeFilter(filter Filter) {
	for i := range fm.filters {
		if fm.filters[i] == filter {
			fm.RemoveFilterByIndex(i)
			return
		}
	}
}

// RemoveFilter remove the filter from this FilterManager
func (fm *filterManager) RemoveFilter(filter ...Filter) {
	for i := range filter {
		fm.removeFilter(filter[i])
	}
}

func (fm *filterManager) inputFilter(data []byte, context Context) []byte {
	for i := len(fm.filters) - 1; i >= 0; i-- {
		data = fm.filters[i].InputFilter(data, context)
	}
	return data
}

func (fm *filterManager) outputFilter(data []byte, context Context) []byte {
	for i := range fm.filters {
		data = fm.filters[i].OutputFilter(data, context)
	}
	return data
}
