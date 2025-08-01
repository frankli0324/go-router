package router

import (
	"errors"

	"github.com/frankli0324/go-router/benchmark/router/radix"
)

type Router[T comparable] struct {
	tree            *radix.Tree[T]
	registeredPaths []string
}

func (r *Router[T]) Handle(path string, f T) {
	r.registeredPaths = append(r.registeredPaths, path)

	optionalPaths := getOptionalPaths(path)
	// if not has optional paths, adds the original
	if len(optionalPaths) == 0 {
		r.tree.Add(path, f)
	} else {
		for _, p := range optionalPaths {
			r.tree.Add(p, f)
		}
	}
}

// Handler makes the router implement the http.Handler interface.
func (r *Router[T]) Handler(path string) (zero T, err error) {
	if r == nil {
		return zero, errors.New("not started")
	}
	if len(path) > 1 && path[len(path)-1] == '/' {
		// pipo handler modified logic, always route without tsr
		path = path[:len(path)-1]
	}
	// Try to search in the wild method tree
	if handler := r.tree.Get(path); handler != zero {
		return handler, nil
	}
	return zero, errors.New("no match")
}

func New[T comparable]() *Router[T] {
	return &Router[T]{tree: radix.New[T]()}
}
