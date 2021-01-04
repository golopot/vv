package vv

import (
	"fmt"
	"strings"
)

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
