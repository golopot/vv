package vv

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// DisallowExtraFields checks for extra fields in the json input. This method
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

func indexPath(i int) string {
	return fmt.Sprintf("[%v]", i)
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

func (o *Object) Float64(path ...string) *Float64Value {
	iv, ok := o.getValue(path)
	r := Float64Value{
		root:     o,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     append(o.path, path...),
	}
	markPath(o, path)
	return &r
}

func (o *Object) JsonNumber(path ...string) *JsonNumberValue {
	iv, ok := o.getValue(path)
	r := JsonNumberValue{
		root:     o,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     append(o.path, path...),
	}
	markPath(o, path)
	return &r
}

func (o *Object) Bool(path ...string) *BoolValue {
	iv, ok := o.getValue(path)
	r := BoolValue{
		root:     o,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     append(o.path, path...),
	}
	markPath(o, path)
	return &r
}

func (o *Object) Slice(path ...string) *SliceValue {
	iv, ok := o.getValue(path)
	v := SliceValue{
		root:     o.root,
		iv:       iv,
		isAbsent: !ok,
		isNull:   ok && iv == nil,
		path:     append(o.path, path...),
	}

	markPath(o, path)
	return &v
}
