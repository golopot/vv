package vv

import "encoding/json"

type JsonNumberValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, json.Number) ValidationError
	default_   *json.Number
}

func (v *JsonNumberValue) Default(s json.Number) *JsonNumberValue {
	v.default_ = &s
	return v
}

func (v *JsonNumberValue) Done() json.Number {
	zero := json.Number("")
	if v.isAbsent || v.isNull {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.root.setMissingError(v.path)
			return zero
		}
	}

	tv, ok := v.iv.(json.Number)
	if !ok {
		v.root.setWrongTypeError(v.path, "number")
		return zero
	}

	for _, fn := range v.validators {
		err := fn(v.path, tv)
		if err != nil {
			v.root.setError(err)
			return zero
		}
	}

	return tv
}

func (v *JsonNumberValue) Pipe(fn func(path []string, value json.Number) ValidationError) *JsonNumberValue {
	v.validators = append(v.validators, fn)
	return v
}
