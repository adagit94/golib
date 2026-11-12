# Overview

Router implemented in Go. Agnostic of specific http implementation.

## Examples
Little wrapper is needed for usage in context of specific http implementation like fasthttp.

### fasthttp
    package server

    import (
      "github.com/adagit94/golib/http"
      "github.com/valyala/fasthttp"
    )

    type reqHandler func(ctx *fasthttp.RequestCtx, params http.PathParams)

    func main() {
      router := goter.CreateRouter[reqHandler]()

      router.Route("a/:b").Get(func(ctx *fasthttp.RequestCtx, params http.PathParams) {
        b := PathParams["b"]
      })

      fasthttp.ListenAndServe("localhost:4000", func(ctx *fasthttp.RequestCtx) {
        handler, params := router.Select(string(ctx.Path()), string(ctx.Method()))
        handler(ctx, params)
      })
    }