package router

import (
	"testing"

	"github.com/fasthttp/router/radix"
	ar "github.com/frankli0324/go-router"
	"github.com/frankli0324/go-router/benchmark/router"
	"github.com/valyala/fasthttp"
)

type testRequests []struct {
	path       string
	nilHandler bool
	route      string
	ps         map[string]string
}

func BenchmarkMatch(b *testing.B) {
	routes := [...]string{
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
	tree := ar.NewRouter[fasthttp.RequestHandler]()
	r := radix.New()
	rr := router.New[*fasthttp.RequestHandler]()
	var zero fasthttp.RequestHandler
	for _, route := range routes {
		tree.Handle(route, func(ctx *fasthttp.RequestCtx) {})
		r.Add(route, func(ctx *fasthttp.RequestCtx) {})
		rr.Handle(route, &zero)
	}

	t := testRequests{
		{"/a", false, "/a", nil},
		{"/", true, "", nil},
		{"/hi", false, "/hi", nil},
		{"/contact", true, "", nil}, // TSR
		{"/co", false, "/co", nil},
		{"/con", true, "", nil},  // key mismatch
		{"/cona", true, "", nil}, // key mismatch
		{"/no", true, "", nil},   // no matching child
		{"/ab", false, "/ab", nil},
		{"/α", false, "/α", nil},
		{"/β", false, "/β", nil},
		{"/hello/test", false, "/hello/test", nil},
		{"/hello/tooth", false, "/hello/tooth", nil},
		{"/hello/testastretta", false, "/hello/{name}", map[string]string{"name": "testastretta"}},
		{"/hello/tes", false, "/hello/{name}", map[string]string{"name": "tes"}},
		{"/hello/test/bye", true, "", nil},
		{"/regex/more_alt/hello", false, "/regex/{path:*}", map[string]string{"path": "more_alt/hello"}},
		{"/regex/small_alt/hello", false, "/regex/{c1:big_alt|alt|small_alt}/{rest:*}", map[string]string{"c1": "small_alt", "rest": "hello"}},
		{"/regex/small_alt/hello", false, "/regex/{c1:big_alt|alt|small_alt}/{rest:*}", map[string]string{"c1": "small_alt", "rest": "hello"}},
		{"/regex/extra_alt/hello", false, "/regex/{c2:(?<named>extra)_alt}/{rest:*}", map[string]string{"c2": "extra_alt", "named": "extra", "rest": "hello"}}, // named group
		{"/wildcard/sub", false, "/wildcard/{rest:*}", map[string]string{"rest": "sub"}},
	}

	b.Run("fasthttp", func(b *testing.B) {
		for i := 0; i < 2*b.N; i++ {
			for _, request := range t {
				r.Get(request.path, nil)
			}
		}
	})
	b.Run("router", func(b *testing.B) {
		for i := 0; i < 2*b.N; i++ {
			for _, request := range t {
				tree.GetParam(request.path, nil)
			}
		}
	})
	b.Run("RouteHandler", func(b *testing.B) {
		for i := 0; i < 2*b.N; i++ {
			for _, request := range t {
				rr.Handler(request.path)
			}
		}
	})
}
