package vv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Object struct {
	iv            interface{}
	root          *Object
	path          []string
	checkedFields map[string]struct{}
	err           ValidationError
}

func New(source []byte) *Object {
	d := json.NewDecoder(bytes.NewReader(source))
	d.UseNumber()

	o := Object{}
	o.checkedFields = map[string]struct{}{}
	o.root = &o
	err := d.Decode(&o.iv)

	if err != nil {
		o.err = JsonParseError{}
	}
	return &o
}

func newObject(iv interface{}, root *Object, path []string) *Object {
	return &Object{
		iv:            iv,
		root:          root,
		path:          path,
		checkedFields: map[string]struct{}{},
	}
}

// ValidationError returns the first error occured
func (o *Object) ValidationError() ValidationError {
	return o.err
}

type ValidationError interface {
	Error() string
	IsValidationError() // A noop method for implementing the interface.
}

type JsonParseError struct{}

func (e JsonParseError) Error() string {
	return "JSON parse error"
}

func (e JsonParseError) IsValidationError() {}

type MissingError struct {
	Path []string
}

func (e MissingError) Error() string {
	return fmt.Sprintf("property `%v` is missing", strings.Join(e.Path, "."))
}

func (e MissingError) IsValidationError() {}

type WrongTypeError struct {
	Path     []string
	Expected string
}

func (e WrongTypeError) Error() string {
	return fmt.Sprintf("property `%v` should be %v", strings.Join(e.Path, "."), e.Expected)
}

func (e WrongTypeError) IsValidationError() {}

type ExtraFieldError struct {
	Path []string
}

func (e ExtraFieldError) Error() string {
	return fmt.Sprintf("extra field `%v` is not allowed", strings.Join(e.Path, "."))
}

func (e ExtraFieldError) IsValidationError() {}

func (v *Object) getValue(path []string) (interface{}, bool) {
	cur := v.iv
	for _, field := range path {
		if obj, ok := cur.(map[string]interface{}); ok {
			if cur, ok = obj[field]; !ok {
				return nil, false
			}
		} else {
			return nil, false
		}
	}

	return cur, true
}

func (o *Object) setMissingError(path []string) {
	if o.err == nil {
		o.err = MissingError{
			Path: path,
		}
	}
}

func (o *Object) setWrongTypeError(path []string, expected string) {
	if o.err == nil {
		o.err = WrongTypeError{
			Path:     path,
			Expected: expected,
		}
	}
}

func (o *Object) setError(err ValidationError) {
	if o.err == nil {
		o.err = err
	}
}

// DisallowExtraFields checks for extra fields in the json input. This function
// must be called **after** all fields are checked.
func (v *Object) DisallowExtraFields() {
	obj, ok := v.iv.(map[string]interface{})
	if !ok {
		return
	}

	for k := range obj {
		if _, ok := v.checkedFields[k]; !ok {
			if v.err == nil {
				v.err = ExtraFieldError{Path: []string{k}}
				return
			}
		}
	}
}

func markPath(v *Object, path []string) {
	if len(path) == 0 {
		return
	}
	key := path[0]
	v.checkedFields[key] = struct{}{}
	for _, name := range path[1:] {
		key += "." + name
		v.checkedFields[key] = struct{}{}
	}
}

type StringValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, string) ValidationError
	default_   *string
}

func (o *Object) String(path ...string) *StringValue {
	iv, ok := o.getValue(path)

	r := StringValue{
		root:     o,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     append(o.path, path...),
	}

	markPath(o, path)

	return &r
}

func (v *StringValue) Default(s string) *StringValue {
	v.default_ = &s
	return v
}

func (v *StringValue) Done() string {
	if v.isAbsent || v.isNull {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.root.setMissingError(v.path)
			return ""
		}
	}

	tv, ok := v.iv.(string)

	if !ok {
		v.root.setWrongTypeError(v.path, "string")
		return ""
	}

	for _, fn := range v.validators {
		err := fn(v.path, tv)
		if err != nil {
			v.root.setError(err)
			return ""
		}
	}

	return tv
}

func (v *StringValue) Pipe(fn func(path []string, value string) ValidationError) *StringValue {
	v.validators = append(v.validators, fn)
	return v
}

type IntValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, int) ValidationError
	default_   *int
}

func (o *Object) Int(path ...string) *IntValue {
	iv, ok := o.getValue(path)
	v := IntValue{
		root:     o.root,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     append(o.path, path...),
	}
	markPath(o, path)
	return &v
}

func (v *IntValue) Default(s int) *IntValue {
	v.default_ = &s
	return v
}

func (v *IntValue) Done() int {
	if v.isAbsent || v.isNull {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.root.setMissingError(v.path)
			return 0
		}
	}

	tv, ok := v.iv.(json.Number)
	if !ok {
		v.root.setWrongTypeError(v.path, "int")
		return 0
	}

	num, err := tv.Int64()

	if err != nil {
		v.root.setWrongTypeError(v.path, "int")
		return 0
	}

	for _, fn := range v.validators {
		err := fn(v.path, int(num))
		if err != nil {
			v.root.setError(err)
			return 0
		}
	}

	return int(num)
}

func (v *IntValue) Pipe(fn func(path []string, value int) ValidationError) *IntValue {
	v.validators = append(v.validators, fn)
	return v
}

func indexPath(i int) string {
	return fmt.Sprintf("[%v]", i)
}

type SliceValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, []interface{}) ValidationError
	default_   *[]interface{}
}

func (o *Object) Slice(path ...string) *SliceValue {
	iv, ok := o.getValue(path)
	v := SliceValue{
		root:     o.root,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     path,
	}

	markPath(o, path)
	return &v
}

func (v *SliceValue) Default(s []interface{}) *SliceValue {
	v.default_ = &s
	return v
}

func (v *SliceValue) Done() []*Object {
	if v.isAbsent || v.isNull {
		if v.default_ != nil {
			list := []*Object{}
			for i, item := range *v.default_ {
				o := newObject(item, v.root, append(v.path, indexPath(i)))
				list = append(list, o)
			}
			return list
		} else {
			v.root.setMissingError(v.path)
			return nil
		}
	}

	tv, ok := v.iv.([]interface{})
	if !ok {
		v.root.setWrongTypeError(v.path, "array")
		return nil
	}

	for _, fn := range v.validators {
		err := fn(v.path, tv)
		if err != nil {
			v.root.setError(err)
			return nil
		}
	}

	result := []*Object{}
	for i, item := range tv {
		any := newObject(item, v.root, append(v.path, indexPath(i)))
		result = append(result, any)
	}

	return result
}

func (v *SliceValue) Pipe(fn func(path []string, value []interface{}) ValidationError) *SliceValue {
	v.validators = append(v.validators, fn)
	return v
}

type Float64Value struct {
	parent   *Object
	path     []string
	default_ *float64
}

type BoolValue struct {
	parent   *Object
	path     []string
	default_ *bool
}
