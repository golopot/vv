package vv

type Float64Value struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, float64) ValidationError
	default_   *float64
}

func (v *Float64Value) Default(w float64) *Float64Value {
	v.default_ = &w
	return v
}

func (v *Float64Value) Done() float64 {
	zero := float64(0)

	if v.isAbsent || v.isNull {
		if v.default_ != nil {
			return *v.default_
		} else {
			v.root.setMissingError(v.path)
			return zero
		}
	}

	tv, ok := v.iv.(float64)

	if !ok {
		v.root.setWrongTypeError(v.path, "float64")
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

func (v *Float64Value) Pipe(fn func(path []string, value float64) ValidationError) *Float64Value {
	v.validators = append(v.validators, fn)
	return v
}
