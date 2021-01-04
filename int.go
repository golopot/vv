package vv

import "encoding/json"

type IntValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, int) ValidationError
	default_   *int
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
