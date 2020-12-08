package vv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Object struct {
	v             map[string]interface{}
	checkedFields map[string]struct{}
	err           ValidationError
}

func New(source []byte) Object {
	d := json.NewDecoder(bytes.NewReader(source))
	d.UseNumber()

	o := Object{
		checkedFields: map[string]struct{}{},
	}
	err := d.Decode(&o.v)

	if err != nil {
		o.err = JsonParseError{}
	}
	return o
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
	iv := v.v
	var result interface{}

	for _, field := range path {
		exists := false
		result, exists = iv[field]
		if !exists {
			return result, false
		}

		deeper := false
		if iv, deeper = result.(map[string]interface{}); !deeper {
			break
		}
	}

	return result, true
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

// CheckExtraFields checks for extra fields in the json input. This function
// must be called **after** all fields are checked.
func (v *Object) CheckExtraFields() {
	for k := range v.v {
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
	parent   *Object
	path     []string
	default_ *string
}

func (o *Object) String(path ...string) *StringValue {
	v := StringValue{
		parent: o,
		path:   path,
	}

	markPath(o, path)

	return &v
}

func (v *StringValue) Default(s string) *StringValue {
	v.default_ = &s
	return v
}

func (v *StringValue) Done() string {
	iv, ok := v.parent.getValue(v.path)
	if !ok {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.parent.setMissingError(v.path)
			return ""
		}
	}

	tv, ok := iv.(string)

	if !ok {
		v.parent.setWrongTypeError(v.path, "string")
		return ""
	}

	return tv
}

type IntValue struct {
	parent   *Object
	path     []string
	default_ *int
}

func (o *Object) Int(path ...string) *IntValue {
	v := IntValue{
		parent: o,
		path:   path,
	}
	markPath(o, path)
	return &v
}

func (v *IntValue) Default(s int) *IntValue {
	v.default_ = &s
	return v
}

func (v *IntValue) Done() int {
	iv, ok := v.parent.getValue(v.path)
	if !ok {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.parent.setMissingError(v.path)
			return 0
		}
	}

	tv, ok := iv.(json.Number)
	if !ok {
		v.parent.setWrongTypeError(v.path, "int")
		return 0
	}

	num, err := tv.Int64()

	if err != nil {
		v.parent.setWrongTypeError(v.path, "int")
		return 0
	}

	return int(num)
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
