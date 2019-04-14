/*
 * Copyright 2012-2019 Li Kexian
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * A toolkit for Golang development
 * https://www.likexian.com/
 */

package xstring

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Version returns package version
func Version() string {
	return "0.3.0"
}

// Author returns package author
func Author() string {
	return "[Li Kexian](https://www.likexian.com/)"
}

// License returns package license
func License() string {
	return "Licensed under the Apache License 2.0"
}

// IsLetter returns if s is an english letter
func IsLetter(s uint8) bool {
	n := (s | 0x20) - 'a'
	return n >= 0 && n < 26
}

// IsLetters returns if s is all english letter
func IsLetters(s string) bool {
	for _, v := range s {
		if !IsLetter(uint8(v)) {
			return false
		}
	}

	return true
}

// IsNumeric returns if s is a number
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// Reverse returns reversed string
func Reverse(s string) string {
	n := len(s)

	runes := make([]rune, n)
	for _, v := range s {
		n--
		runes[n] = v
	}

	return string(runes[n:])
}

// ToString convert v to string
func ToString(v interface{}) string {
	switch vv := v.(type) {
	case []byte:
		return string(vv)
	case string:
		return vv
	case bool:
		return strconv.FormatBool(vv)
	case int:
		return strconv.FormatInt(int64(vv), 10)
	case int8:
		return strconv.FormatInt(int64(vv), 10)
	case int16:
		return strconv.FormatInt(int64(vv), 10)
	case int32:
		return strconv.FormatInt(int64(vv), 10)
	case int64:
		return strconv.FormatInt(int64(vv), 10)
	case uint:
		return strconv.FormatUint(uint64(vv), 10)
	case uint8:
		return strconv.FormatUint(uint64(vv), 10)
	case uint16:
		return strconv.FormatUint(uint64(vv), 10)
	case uint32:
		return strconv.FormatUint(uint64(vv), 10)
	case uint64:
		return strconv.FormatUint(uint64(vv), 10)
	case float32:
		return strconv.FormatFloat(float64(vv), 'f', 2, 64)
	case float64:
		return strconv.FormatFloat(float64(vv), 'f', 2, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Join concatenates the elements and returns string
func Join(v interface{}, sep string) string {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Ptr || vv.Kind() == reflect.Interface {
		if vv.IsNil() {
			return ""
		}
		vv = vv.Elem()
	}

	switch vv.Kind() {
	case reflect.Slice, reflect.Array:
		as := []string{}
		for i := 0; i < vv.Len(); i++ {
			as = append(as, ToString(vv.Index(i)))
		}
		return strings.Join(as, sep)
	default:
		return ToString(v)
	}
}

// Expand replaces {var} of string s based on the value map m
// For example, Expand("i am {name}", map[string]interface{}{"name": "Li Kexian"})
func Expand(s string, m map[string]interface{}) string {
	var i, j int
	var buf []byte

	for {
		i = LastInIndex(s, "{")
		if i < 0 {
			break
		}
		j = strings.Index(s[i+1:], "}")
		if j <= 0 {
			break
		}
		buf = append(buf, s[:i]...)
		key := s[i+1 : i+1+j]
		if v, ok := m[key]; ok {
			buf = append(buf, fmt.Sprint(v)...)
		} else {
			buf = append(buf, []byte(fmt.Sprintf("%%!%s(MISSING)", key))...)
		}
		s = s[i+1+j+1:]
	}

	buf = append(buf, s...)
	s = string(buf)

	return s
}

// LastInIndex find last position at first index
//   for example, LastInIndex("{{{{{{{{{{name}", "{")
//                                      ↑
func LastInIndex(s, f string) int {
	i := strings.Index(s, f)
	if i < 0 {
		return i
	}

	t := s[i+1:]
	for j := 0; j < len(t); j++ {
		if t[j] != f[0] {
			return j + i
		}
	}

	return i
}
