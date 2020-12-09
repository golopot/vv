package main

import (
	"fmt"

	"github.com/golopot/vv"
)

// func mustStartsWithA(w vv.StringValue) vv.ValidationError {
// 	s := w.Done()

// 	if s[0] != 'A' {
// 		return vv.NewCustomValidationError(w, "%v must starts with A.")
// 	}

// 	return nil
// }

func errorPrinter(err vv.ValidationError) string {
	switch e := err.(type) {
	case vv.JsonParseError:
		return "invalid json"
	case vv.WrongTypeError:
		return fmt.Sprintf("%v should be %v", e.Path, e.Expected)
	}
	return err.Error()
}

func main() {
	input := []byte(`
{
	"foo": 123,
	"bar": "qqq",
	"a": {
		"b": "ewq"
	}
}`)

	v := vv.New(input)
	foo := v.Int("foo").Done()
	bar := v.String("bar").Done()

	// for _, w := range v.Slice("goo").Done() {
	// 	q := w.Int().Done()
	// }
	// goo := v.String("goo").Pipe(mustStartsWithA).Done()

	if err := v.ValidationError(); err != nil {
		fmt.Println("error:", errorPrinter(err))
		return
	}

	fmt.Println(foo, bar)
}
