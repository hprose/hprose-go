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
 * hprose/struct_encoder.go                               *
 *                                                        *
 * hprose struct encoder for Go.                          *
 *                                                        *
 * LastModified: Aug 29, 2015                             *
 * Author: Ma Bingyao <andot@hprose.com>                  *
 *                                                        *
\**********************************************************/

package io

import (
	"bytes"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

type fieldCache struct {
	Name   string
	Alias  string
	Index  []int
	Offset uintptr
	Type   uintptr
	Kind   reflect.Kind
}

type structCache struct {
	Alias  string
	Tag    string
	Fields []*fieldCache
	Data   []byte
}

var structTypeCache = map[uintptr]*structCache{}
var structTypeCacheLocker = sync.RWMutex{}

var structTypes = map[string]reflect.Type{}
var structTypesLocker = sync.RWMutex{}

func getFieldAlias(f *reflect.StructField, tag string) (alias string) {
	fname := f.Name
	if fname != "" && 'A' <= fname[0] && fname[0] < 'Z' {
		if tag != "" && f.Tag != "" {
			alias = strings.SplitN(f.Tag.Get(tag), ",", 2)[0]
			alias = strings.TrimSpace(strings.SplitN(alias, ">", 2)[0])
			if alias == "-" {
				return ""
			}
		}
		if alias == "" {
			alias = string(fname[0]-'A'+'a') + fname[1:]
		}
	}
	return alias
}

func getSubFields(t reflect.Type, tag string, offset uintptr, index []int) []*fieldCache {
	subFields := getFields(t, tag)
	for _, subField := range subFields {
		subField.Offset += offset
		subField.Index = append(index, subField.Index...)
	}
	return subFields
}

func getFields(t reflect.Type, tag string) []*fieldCache {
	n := t.NumField()
	fields := make([]*fieldCache, 0, n)
	for i := 0; i < n; i++ {
		f := t.Field(i)
		ft := f.Type
		fkind := ft.Kind()
		if fkind == reflect.Chan ||
			fkind == reflect.Func ||
			fkind == reflect.UnsafePointer {
			continue
		}
		if f.Anonymous {
			if fkind == reflect.Struct {
				subFields := getSubFields(ft, tag, f.Offset, f.Index)
				fields = append(fields, subFields...)
				continue
			}
		}
		alias := getFieldAlias(&f, tag)
		if alias == "" {
			continue
		}
		field := fieldCache{}
		field.Name = f.Name
		field.Alias = alias
		field.Type = (*emptyInterface)(unsafe.Pointer(&ft)).ptr
		field.Kind = fkind
		field.Offset = f.Offset
		field.Index = f.Index
		fields = append(fields, &field)
	}
	return fields
}

func initStructCacheData(cache *structCache) {
	s := &bytes.Buffer{}
	fields := cache.Fields
	count := len(fields)
	s.WriteByte(TagClass)
	var buf [20]byte
	s.Write(getIntBytes(buf[:], int64(utf16Length(cache.Alias))))
	s.WriteByte(TagQuote)
	s.WriteString(cache.Alias)
	s.WriteByte(TagQuote)
	if count > 0 {
		s.Write(getIntBytes(buf[:], int64(count)))
	}
	s.WriteByte(TagOpenbrace)
	for _, field := range fields {
		s.WriteByte(TagString)
		s.Write(getIntBytes(buf[:], int64(utf16Length(field.Alias))))
		s.WriteByte(TagQuote)
		s.WriteString(field.Alias)
		s.WriteByte(TagQuote)
	}
	s.WriteByte(TagClosebrace)
	cache.Data = s.Bytes()
}

func getStructCache(structType reflect.Type) *structCache {
	typ := (*emptyInterface)(unsafe.Pointer(&structType)).ptr
	structTypeCacheLocker.RLock()
	cache, ok := structTypeCache[typ]
	if !ok {
		structTypeCacheLocker.RUnlock()
		structTypeCacheLocker.Lock()
		cache, ok = structTypeCache[typ]
		if !ok {
			cache = &structCache{}
			cache.Alias = structType.Name()
			cache.Fields = getFields(structType, "")
			initStructCacheData(cache)
			structTypeCache[typ] = cache
		}
		structTypeCacheLocker.Unlock()
	} else {
		structTypeCacheLocker.RUnlock()
	}
	return cache
}

// Register structType with alias & tag.
func Register(structType reflect.Type, alias string, tag ...string) {
	if structType.Kind() == reflect.Ptr {
		structType = structType.Elem()
	}
	if structType.Kind() != reflect.Struct {
		panic("invalid type: " + structType.String())
	}
	structTypesLocker.Lock()
	structTypes[alias] = structType
	structTypesLocker.Unlock()

	structTypeCacheLocker.Lock()
	cache := &structCache{Alias: alias}
	if len(tag) == 1 {
		cache.Tag = tag[0]
	}
	cache.Fields = getFields(structType, cache.Tag)
	initStructCacheData(cache)
	structTypeCache[(*emptyInterface)(unsafe.Pointer(&structType)).ptr] = cache
	structTypeCacheLocker.Unlock()
}

// GetStructType by alias.
func GetStructType(alias string) (structType reflect.Type) {
	structTypesLocker.RLock()
	structType = structTypes[alias]
	structTypesLocker.RUnlock()
	return structType
}

// GetAlias of structType
func GetAlias(structType reflect.Type) string {
	return getStructCache(structType).Alias
}

// GetTag by structType.
func GetTag(structType reflect.Type) string {
	return getStructCache(structType).Tag
}
