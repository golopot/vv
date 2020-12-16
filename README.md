# VV

[![GoDoc](https://pkg.go.dev/badge/github.com/golopot/vv)](https://pkg.go.dev/github.com/golopot/vv)

A json validation library.

- Get json value and validate json at the same time.
- Concise usage.
- Practical error handling.
- Customizable error message.

## Features

- Optional fields and default value
- Disallow extra object fields
- User defined validators
- Customizable error message

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
	"a": 1,
	"b": "sss"
}
`))

    a := v.Int("a").Done()
	z := v.String("z").Default("some-string").Done() // `z` is optional
	wrong := v.Int("wrong").Done() // fails validation

	fmt.Println(a, z, wrong) // 1 "some-string" 0

	// v.ValidationError() stores the first error occured
	if err := v.ValidationError(); err != nil {
        fmt.Println("error:", err)
        // error: property `wrong` should be int
		return
	}
}
```

## License

MIT
