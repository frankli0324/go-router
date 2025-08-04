package router

import (
	"testing"

	"github.com/fasthttp/router/radix"
	ar "github.com/frankli0324/go-router"
	"github.com/frankli0324/go-router/benchmark/router"
	"github.com/valyala/fasthttp"
)

func execute(b *testing.B, routes []string, t []string) {
	tree := ar.NewRouter[fasthttp.RequestHandler]()
	r := radix.New()
	rr := router.New[*fasthttp.RequestHandler]()
	var zero fasthttp.RequestHandler
	for _, route := range routes {
		tree.Set(route, func(ctx *fasthttp.RequestCtx) {})
		r.Add(route, func(ctx *fasthttp.RequestCtx) {})
		rr.Handle(route, &zero)
	}

	b.Run("RouteHandler", func(b *testing.B) {
		for i := 0; i < 2*b.N; i++ {
			for _, request := range t {
				rr.Handler(request)
			}
		}
	})
	b.Run("router", func(b *testing.B) {
		for i := 0; i < 2*b.N; i++ {
			for _, request := range t {
				tree.GetParam(request, nil)
			}
		}
	})
	b.Run("fasthttp", func(b *testing.B) {
		for i := 0; i < 2*b.N; i++ {
			for _, request := range t {
				r.Get(request, nil)
			}
		}
	})
}

func Benchmark(b *testing.B) {
	routes := []string{
		"/hi",
		"/contact/",
		"/co",
		"/c",
		"/a",
		"/ab",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/α",
		"/β",
		"/hello/test",
		"/hello/tooth",
		"/hello/{name}",
		"/regex/{c1:big_alt|alt|small_alt}/{rest:*}",
		"/regex/{c2:(?<named>extra)_alt}/{rest:*}",
		"/regex/{path:*}",
		"/wildcard/sub/{rest:*}",
		"/wildcard/{rest:*}",
	}

	t := []string{
		"/a",
		"/",
		"/hi",
		"/contact",
		"/co",
		"/con",
		"/cona",
		"/no",
		"/ab",
		"/α",
		"/β",
		"/hello/test",
		"/hello/tooth",
		"/hello/testastretta",
		"/hello/tes",
		"/hello/test/bye",
		"/regex/more_alt/hello",
		"/regex/small_alt/hello",
		"/regex/small_alt/hello",
		"/regex/extra_alt/hello",
		"/wildcard/sub",
	}
	execute(b, routes[:], t)
}
