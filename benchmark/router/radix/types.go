package radix

import (
	"context"
	"regexp"
)

type VType[T any] func(context.Context, T) error

type nodeType uint8

type nodeWildcard[T comparable] struct {
	path     string
	paramKey string
	handler  T
}

type node[T comparable] struct {
	nType nodeType

	path         string
	handler      T
	hasWildChild bool
	children     []*node[T]
	wildcard     *nodeWildcard[T]

	paramKeys  []string
	paramRegex *regexp.Regexp
}

type wildPath struct {
	path  string
	keys  []string
	start int
	end   int
	pType nodeType

	pattern string
	regex   *regexp.Regexp
}

// Tree is a routes storage
type Tree[T comparable] struct {
	root *node[T]

	// If enabled, the node handler could be updated
	Mutable bool
}
