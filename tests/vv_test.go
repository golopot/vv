package tests

import (
	"fmt"
	"testing"

	"github.com/golopot/vv"

	"github.com/stretchr/testify/assert"
)

func noop(...interface{}) {}

func TestPass(t *testing.T) {

	v := vv.New([]byte(`
	{
		"a": "ssss",
		"b": 123
	}
`))

	a := v.String("a").Done()
	b := v.Int("b").Done()

	if err := v.ValidationError(); err != nil {
		t.Errorf("expected validation pass, but got: %v", err)
		return
	}

	assert.Equal(t, a, "ssss")
	assert.Equal(t, b, 123)
}

func TestJsonParseError(t *testing.T) {

	v := vv.New([]byte(`
	{
		"a": 1,
	}
`))

	a := v.String("a").Done()
	noop(a)

	assert.Equal(
		t,
		vv.JsonParseError{},
		v.ValidationError(),
	)
}

func TestMissing(t *testing.T) {

	v := vv.New([]byte(`
	{
		"a": "ssss"
	}
`))

	a := v.String("a").Done()
	b := v.Int("b").Done()
	noop(a, b)

	assert.Equal(
		t,
		vv.MissingError{
			Path: []string{"b"},
		},
		v.ValidationError(),
	)
}

func TestWrongType(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": "ssss"
	}
`))

	a := v.Int("a").Done()
	noop(a)

	assert.Equal(
		t,
		vv.WrongTypeError{
			Path:     []string{"a"},
			Expected: "int",
		},
		v.ValidationError(),
	)
}

func TestOptionalPass(t *testing.T) {
	v := vv.New([]byte(`
	{

	}
`))

	a := v.Int("a").Default(123).Done()

	assert.Equal(t, nil, v.ValidationError())
	assert.Equal(t, a, 123)
}

func TestOptionalFail(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": "s"
	}
`))

	a := v.Int("a").Default(123).Done()
	noop(a)

	assert.Equal(
		t,
		vv.WrongTypeError{
			Path:     []string{"a"},
			Expected: "int",
		},
		v.ValidationError(),
	)
}

func TestOptionalNull(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": null
	}
`))

	a := v.Int("a").Default(123).Done()
	noop(a)
	assert.Equal(t, nil, v.ValidationError())
	assert.Equal(t, 123, a)
}

func TestExtraPass(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": "s"
	}
`))

	a := v.String("a").Done()
	noop(a)

	v.DisallowExtraFields()

	assert.Equal(
		t,
		nil,
		v.ValidationError(),
	)
}

func TestExtraFail(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": "s",
		"b": 5
	}
`))

	a := v.String("a").Done()
	noop(a)

	v.DisallowExtraFields()

	assert.Equal(
		t,
		vv.ExtraFieldError{
			Path: []string{"b"},
		},
		v.ValidationError(),
	)
}

type LengthIsNotFourError struct {
	Path []string
}

func (LengthIsNotFourError) IsValidationError() {}
func (l LengthIsNotFourError) Error() string {
	return fmt.Sprintf("Length of %v must be 4.", l.Path)
}

func LengthIsFour(path []string, v string) vv.ValidationError {
	if len(v) != 4 {
		return LengthIsNotFourError{Path: path}
	}
	return nil
}

func TestPipeFail(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": "12345"
	}
`))

	a := v.String("a").Pipe(LengthIsFour).Done()
	noop(a)

	assert.Equal(
		t,
		LengthIsNotFourError{[]string{"a"}},
		v.ValidationError(),
	)
}

func TestPipePass(t *testing.T) {
	v := vv.New([]byte(`
	{
		"a": "1234"
	}
`))

	a := v.String("a").Pipe(LengthIsFour).Done()
	noop(a)

	assert.Equal(t, nil, v.ValidationError())
}
