package vv

type BoolValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, bool) ValidationError
	default_   *bool
}

func (v *BoolValue) Default(w bool) *BoolValue {
	v.default_ = &w
	return v
}

func (v *BoolValue) Done() bool {
	zero := false

	if v.isAbsent || v.isNull {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.root.setMissingError(v.path)
			return zero
		}
	}

	tv, ok := v.iv.(bool)

	if !ok {
		v.root.setWrongTypeError(v.path, "boolean")
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

func (v *BoolValue) Pipe(fn func(path []string, value bool) ValidationError) *BoolValue {
	v.validators = append(v.validators, fn)
	return v
}
