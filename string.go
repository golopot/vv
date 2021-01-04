package vv

type StringValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, string) ValidationError
	default_   *string
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
