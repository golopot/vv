# VV

A JSON validation library.

## Example

```go
package server

import (
	"github.com/golopot/vv"
)

func someHandler(ctx *fasthttp.RequestCtx) {
    v := vv.New(ctx.Request.Body())

    foo := v.Prop("foo").Int()
    gooo := v.Prop("gooo").Float64()
    moooo := v.Prop("moooo").String()
    nooo := v.Prop("nooo").String(vv.OptionalString("some-default-value"))

    if err := v.ValidationError(); err != nil {
        sendError(ctx, err.Error(), 400)
        return
    }

    ...
}

```
