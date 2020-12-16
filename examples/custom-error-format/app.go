package main

import (
	"fmt"
	"strings"

	"github.com/golopot/vv"
)

func formatValidationError(err vv.ValidationError) string {
	formatPath := func(path []string) string {
		return strings.Join(path, ".")
	}

	switch e := err.(type) {
	case vv.JsonParseError:
		return "invalid json"
	case vv.MissingError:
		return fmt.Sprintf("%v is missing.", formatPath(e.Path))
	case vv.WrongTypeError:
		return fmt.Sprintf("%v should be of type %v.", formatPath(e.Path), e.Expected)
	case vv.ExtraFieldError:
		return fmt.Sprintf("Field %v is not allowed.", formatPath(e.Path))
	default:
		// fallback to default
		return err.Error()
	}
}

func main() {
	v := vv.New([]byte(`
	{
		"a": 1
	}
`))

	a := v.String("a")
	fmt.Println(a)

	if err := v.ValidationError(); err != nil {
		fmt.Println(formatValidationError(err))
		return
	}
}
