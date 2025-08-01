package router

import (
	"fmt"
)

func ExampleNewRouter() {
	r := NewRouter[any]()
	r.Set("/api/v1/hello", 1)
	r.Set("/api/v1/world", 2)
	fmt.Println(r.Get("/api/v1/hello") == 1)
}

func ExampleRouter_Get() {
	r := NewRouter[any]()
	r.Set("/", "notmatch")
	r.Set("/1", 1)
	r.Set("/sub/test", "literals have highest prio")
	r.Set("/sub/{a}", "2")
	r.Set("/{anything:*}", 3.0)
	fmt.Println(r.Get("/sub/test") == "literals have highest prio")
	fmt.Println(r.Get("/sub/other") == "2")
	fmt.Println(r.Get("/sub/other/") == 3.0) // because it doesn't match /sub/{a} due to an extra trailing slash
	fmt.Println(r.Get("/") == 3.0)           // "anything" could be empty, same as fasthttp
}

func ExampleRouter_GetParam() {
	r := NewRouter[any]()
	r.Set("/", 1)
	r.Set("/{a}", "2")
	param := make(map[string]string)
	fmt.Println(r.GetParam("/test", param) == "2")
	fmt.Println(param["a"] == "test")
}
