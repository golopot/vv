# VV

A JSON validation library.

## Example

```go
package main

import (
	"fmt"

	"github.com/golopot/vv"
)

func main() {
	v := vv.New([]byte(`
{
	"foo": 1,
	"goo": "aaa",
	"hoo": "aaa"
}
`))

	foo := v.Int("foo").Done()
	goo := v.String("goo").Done()
	noo := v.String("noo").Default("some-string").Done() // `noo` is optional with a default value
	hoo := v.Int("hoo").Done()

	fmt.Println(foo, goo, noo, hoo) // 1 aaa some-string

	// v.ValidationError() gives the result of validation
	if err := v.ValidationError(); err != nil {
        fmt.Println("error:", err)
        // error: property `hoo` should be int
		return
	}
}
```

## License

MIT
