package vv

type SliceValue struct {
	root       *Object
	iv         interface{}
	isAbsent   bool
	isNull     bool
	path       []string
	validators []func([]string, []interface{}) ValidationError
	default_   *[]interface{}
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
