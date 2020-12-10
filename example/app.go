package main

import (
	"fmt"

	"github.com/golopot/vv"
	"github.com/kr/pretty"
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
	v := vv.New([]byte(`
	{
		"a": [1, 2, 3]
	}
`))

	a := v.Slice("a").Done()
	nums := []int{}
	for _, w := range a {
		u := w.Int().Done()
		nums = append(nums, u)
	}

	pretty.Println(a[0])
}
