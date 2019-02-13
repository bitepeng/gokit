/*
 * Go module for JSON parsing
 * https://www.likexian.com/
 *
 * Copyright 2012-2019, Li Kexian
 * Released under the Apache License, Version 2.0
 *
 */

package simplejson


import (
    "os"
    "io"
    "io/ioutil"
    "bytes"
    "errors"
    "encoding/json"
    "reflect"
    "strings"
    "strconv"
    "time"
)


// storing json data
type Json struct {
    data        interface{}
    escapeHtml  bool
}


// returns package version
func Version() string {
    return "0.9.0"
}


// returns package author
func Author() string {
    return "[Li Kexian](https://www.likexian.com/)"
}


// returns package license
func License() string {
    return "Apache License, Version 2.0"
}


// returns a pointer to a new Json object
//   data_json := New()
//   data_json := New(type Data struct{data string}{"zzz"})
//   data_json := New(map[string]interface{}{"iam": "Li Kexian"})
func New(args ...interface{}) (*Json) {
    switch len(args) {
        case 1:
            return &Json {
                data: args[0],
            }
        default:
            return &Json {
                data: make(map[string]interface{}),
            }
    }
}


// set html escape for escaping of <, >, and & in JSON strings
func (j *Json) SetHtmlEscape(escape bool) {
    j.escapeHtml = escape
}


// loads data from a file, returns a json object
func Load(file string) (j *Json, err error) {
    data, err := ioutil.ReadFile(file)
    if err != nil {
        return
    }

    text := string(data)
    j, err = Loads(text)

    return
}


// dumps json object to a file
func (j *Json) Dump(file string) (bytes int, err error) {
    result, err := j.PrettyDumps()
    if err != nil {
        return
    }

    fd, err := os.OpenFile(file, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0644)
    if err != nil {
        return
    }

    bytes, err = io.WriteString(fd, result)
    fd.Close()

    return
}


// unmarshal json from string, returns json object
func Loads(text string) (j *Json, err error) {
    j = new(Json)

    dec := json.NewDecoder(bytes.NewBuffer([]byte(text)))
    dec.UseNumber()
    err = dec.Decode(&j.data)

    return
}


// marshal json object to string
func (j *Json) Dumps() (result string, err error) {
    return j.doDumps("")
}


// marshal json object to string, with identation
func (j *Json) PrettyDumps() (result string, err error) {
    return j.doDumps(strings.Repeat(" ", 4))
}


// do marshal json to string
func (j *Json) doDumps(indent string) (result string, err error) {
    var buf bytes.Buffer

    enc := json.NewEncoder(&buf)
    enc.SetEscapeHTML(j.escapeHtml)
    enc.SetIndent("", indent)
    err = enc.Encode(j.data)
    if err != nil {
        return
    }

    result = buf.String()
    result = strings.TrimSpace(result)

    return
}


// set key-value to json object, dot(.) separated key is supported
//   json.Set("status", 1)
//   json.Set("status.code", 1)
//   ! NOT SUPPORTED json.Set("result.intlist.3", 666)
func (j *Json) Set(key string, value interface{}) {
    key = strings.TrimSpace(key)
    if key == "" {
        j.data = value
        return
    }

    result, err := j.Map()
    if err != nil {
        return
    }

    keys := strings.Split(key, ".")
    for i:=0; i<len(keys)-1; i++  {
        v := strings.TrimSpace(keys[i])
        if v != "" {
            if _, ok := result[v]; !ok {
                result[v] = make(map[string]interface{})
            }
            result = result[v].(map[string]interface{})
        }
    }

    result[keys[len(keys) - 1]] = value
}


// delete key-value from json object, dot(.) separated key is supported
//   json.Del("status")
//   json.Del("status.code")
//   ! NOT SUPPORTED json.Del("result.intlist.3")
func (j *Json) Del(key string) {
    result, err := j.Map()
    if err != nil {
        return
    }

    var ok bool
    keys := strings.Split(key, ".")
    for i:=0; i<len(keys)-1; i++  {
        v := strings.TrimSpace(keys[i])
        if v != "" {
            if _, ok = result[v]; !ok {
                return
            }
            result, ok = result[v].(map[string]interface{})
            if !ok {
                return
            }
        }
    }

    if _, ok := result[keys[len(keys) - 1]]; ok {
        delete(result, keys[len(keys) - 1])
    }
}


// check json object has key, dot(.) separated key is supported
//   json.Has("status")
//   json.Has("status.code")
//   json.Has("result.intlist.3")
func (j *Json) Has(key string) (bool) {
    result := j

    keys := strings.Split(key, ".")
    for i:=0; i<len(keys); i++  {
        v := strings.TrimSpace(keys[i])
        if v != "" {
            tmp, err := result.Map()
            if err == nil {
                if _, ok := tmp[v]; !ok {
                    return false
                }
                if i == len(keys) - 1 {
                    return true
                }
                result = result.Get(v)
            } else {
                tmp, err := result.Array()
                if err == nil {
                    n, err := strconv.Atoi(v)
                    if err != nil {
                        return false
                    }
                    if n >= len(tmp) {
                        return false
                    }
                    if i == len(keys) - 1 {
                        return true
                    }
                    result = result.GetN(n)
                } else {
                    return false
                }
            }
        }
    }

    return false
}


// returns the pointer to json object by key, dot(.) separated key is supported
//   json.Get("status").Int()
//   json.Get("status.code").Int()
//   json.Get("result.intlist.3").Int()
func (j *Json) Get(key string) (*Json) {
    result := j

    for _, v := range strings.Split(key, ".") {
        v = strings.TrimSpace(v)
        if v != "" {
            tmp, err := result.Map()
            if err == nil {
                if _, ok := tmp[v]; ok {
                    result = &Json{tmp[v], j.escapeHtml}
                } else {
                    return &Json{nil, j.escapeHtml}
                }
            } else {
                _, err := result.Array()
                if err == nil {
                    i, err := strconv.Atoi(v)
                    if err != nil {
                        return &Json{nil, j.escapeHtml}
                    }
                    result = result.GetN(i)
                } else {
                    return &Json{nil, j.escapeHtml}
                }
            }
        }
    }

    return result
}


// returns a pointer to the index of json object
//   json.Get("int_list").GetN(1).Int()
func (j *Json) GetN(i int) (*Json) {
    data, err := j.Array()
    if err == nil {
        if len(data) > i {
            return &Json{data[i], j.escapeHtml}
        }
    }

    return &Json{nil, j.escapeHtml}
}


// return len of json object, -1 if type invalid or error
func (j *Json) Len() (int) {
    if v, err := j.Map(); err == nil {
        return len(v)
    }

    if v, err := j.Array(); err == nil {
        return len(v)
    }

    if v, err := j.String(); err == nil {
        return len(v)
    }

    return -1
}


// return json object is a map
func (j *Json) IsMap() (bool) {
    switch j.data.(type) {
        case map[string]interface{}:
            return true
        default:
            return false
    }
}


// return json object is an array
func (j *Json) IsArray() (bool) {
    switch j.data.(type) {
        case []interface{}:
            return true
        default:
            return false
    }
}


// returns as map from json object
func (j *Json) Map() (result map[string]interface{}, err error) {
    result, ok := (j.data).(map[string]interface{})
    if !ok {
        err = errors.New("assert to map failed")
    }
    return
}


// returns as array from json object
func (j *Json) Array() (result []interface{}, err error) {
    result, ok := (j.data).([]interface{})
    if !ok {
        err = errors.New("assert to array failed")
    }
    return
}


// returns as bool from json object
func (j *Json) Bool() (result bool, err error) {
    result, ok := (j.data).(bool)
    if !ok {
        err = errors.New("assert to bool failed")
    }
    return
}


// returns as string from json object
func (j *Json) String() (result string, err error) {
    result, ok := (j.data).(string)
    if !ok {
        err = errors.New("assert to string failed")
    }
    return
}


// returns as string array from json object
func (j *Json) StringArray() (result []string, err error) {
    data, err := j.Array()
    if err != nil {
        return
    }

    for _, v := range data {
        if v == nil {
            result = append(result, "")
        } else {
            r, ok := v.(string)
            if !ok {
                err = errors.New("assert to []string failed")
                return
            }
            result = append(result, r)
        }
    }

    return
}


// returns as time.Time from json object
// optional args is to set the time string parsing format, time.RFC3339 by default
// if the time is of int, optional args must not set
//   json.Time()
//   json.Time("2006-01-02 15:04:05")
func (j *Json) Time(args ...string) (result time.Time, err error) {
    switch j.data.(type) {
        case string:
            if len(args) > 1 {
                return result, errors.New("Too many arguments")
            }
            r, e := j.String()
            if e != nil {
                return result, e
            }
            format := time.RFC3339
            if len(args) == 1 && strings.TrimSpace(args[0]) != "" {
                format = strings.TrimSpace(args[0])
            }
            return time.ParseInLocation(format, r, time.Local)
        default:
            if len(args) > 0 {
                return result, errors.New("Too many arguments")
            }
            r, e := j.Int64()
            if e != nil {
                return result, e
            }
            return time.Unix(r, 0), nil
    }
}


// returns as float64 from json object
func (j *Json) Float64() (result float64, err error) {
    switch j.data.(type) {
        case json.Number:
            return j.data.(json.Number).Float64()
        case float32, float64:
            return reflect.ValueOf(j.data).Float(), nil
        case int, int8, int16, int32, int64:
            return float64(reflect.ValueOf(j.data).Int()), nil
        case uint, uint8, uint16, uint32, uint64:
            return float64(reflect.ValueOf(j.data).Uint()), nil
        default:
            return 0, errors.New("invalid value type")
    }
}


// returns as int from json object
func (j *Json) Int() (result int, err error) {
    switch j.data.(type) {
        case json.Number:
            r, err := j.data.(json.Number).Int64()
            return int(r), err
        case float32, float64:
            return int(reflect.ValueOf(j.data).Float()), nil
        case int, int8, int16, int32, int64:
            return int(reflect.ValueOf(j.data).Int()), nil
        case uint, uint8, uint16, uint32, uint64:
            return int(reflect.ValueOf(j.data).Uint()), nil
        default:
            return 0, errors.New("invalid value type")
    }
}


// returns as int64 from json object
func (j *Json) Int64() (result int64, err error) {
    switch j.data.(type) {
        case json.Number:
            return j.data.(json.Number).Int64()
        case float32, float64:
            return int64(reflect.ValueOf(j.data).Float()), nil
        case int, int8, int16, int32, int64:
            return reflect.ValueOf(j.data).Int(), nil
        case uint, uint8, uint16, uint32, uint64:
            return int64(reflect.ValueOf(j.data).Uint()), nil
        default:
            return 0, errors.New("invalid value type")
    }
}


// returns as uint64 from json object
func (j *Json) Uint64() (result uint64, err error) {
    switch j.data.(type) {
        case json.Number:
            return strconv.ParseUint(j.data.(json.Number).String(), 10, 64)
        case float32, float64:
            return uint64(reflect.ValueOf(j.data).Float()), nil
        case int, int8, int16, int32, int64:
            return uint64(reflect.ValueOf(j.data).Int()), nil
        case uint, uint8, uint16, uint32, uint64:
            return reflect.ValueOf(j.data).Uint(), nil
        default:
            return 0, errors.New("invalid value type")
    }
}


// returns as bool from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustBool(args ...bool) (bool) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.Bool()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}


// returns as string from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustString(args ...string) (string) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.String()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}


// returns as string from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustStringArray(args ...[]string) ([]string) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.StringArray()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}


// returns as time.Time from json object
// if error return default(if set) or panic
//   json.Time()                                                 // No format,  No default
//   json.Time("2006-01-02 15:04:05")                            // Has format, No default
//   json.Time(time.Unix(1548907870, 0))                         // No format,  Has default
//   json.Time("2006-01-02 15:04:05", time.Unix(1548907870, 0))  // Has format, Has default
func (j *Json) MustTime(args ...interface{}) (time.Time) {
    if len(args) > 2 {
        panic("Too many arguments")
    }

    format := ""
    defset := false
    var defbak time.Time

    for i:=0; i<len(args); i++ {
        switch args[i].(type) {
            case string:
                format = args[i].(string)
            case time.Time:
                defbak = args[i].(time.Time)
                defset = true
            default:
                panic("Invalid argument type")
        }
    }

    var r time.Time
    var err error
    if format != "" {
        r, err = j.Time(format)
    } else {
        r, err = j.Time()
    }

    if err == nil {
        return r
    }

    if defset {
        return defbak
    } else {
        panic(err)
    }
}


// returns as float64 from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustFloat64(args ...float64) (float64) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.Float64()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}


// returns as int from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustInt(args ...int) (int) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.Int()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}


// returns as int64 from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustInt64(args ...int64) (int64) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.Int64()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}


// returns as uint64 from json object with optional default value
// if error return default(if set) or panic
func (j *Json) MustUint64(args ...uint64) (uint64) {
    if len(args) > 1 {
        panic("Too many arguments")
    }

    r, err := j.Uint64()
    if err == nil {
        return r
    }

    if len(args) == 1 {
        return args[0]
    } else {
        panic(err)
    }
}
