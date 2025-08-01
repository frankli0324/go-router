package radix

import (
	"errors"
	"strings"
)

const (
	root nodeType = iota
	static
	param
	wildcard
)

// New returns an empty routes storage
func New[T comparable]() *Tree[T] {
	return &Tree[T]{
		root: &node[T]{
			nType: root,
		},
	}
}

// Add adds a node with the given handle to the path.
//
// WARNING: Not concurrency-safe!
func (t *Tree[T]) Add(path string, handler T) {
	var zero T
	if !strings.HasPrefix(path, "/") {
		panicf("path must begin with '/' in path '%s'", path)
	} else if handler == zero {
		panic("zero handler")
	}

	fullPath := path

	i := longestCommonPrefix(path, t.root.path)
	if i > 0 {
		if len(t.root.path) > i {
			t.root.split(i)
		}

		path = path[i:]
	}

	n, err := t.root.add(path, fullPath, handler)
	if err != nil {
		var radixErr radixError

		if errors.As(err, &radixErr) && t.Mutable {
			switch radixErr.msg {
			case errSetHandler:
				n.handler = handler
				return
			case errSetWildcardHandler:
				n.wildcard.handler = handler
				return
			}
		}

		panic(err)
	}

	if len(t.root.path) == 0 {
		t.root = t.root.children[0]
		t.root.nType = root
	}

	// Reorder the nodes
	t.root.sort()
}

// Get returns the handle registered with the given path (key). The values of
// param/wildcard are saved as ctx.UserValue.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (t *Tree[T]) Get(path string) (zero T) {
	return t.getWithParams(path, map[string]string{})
}

func (t *Tree[T]) getWithParams(path string, params map[string]string) (zero T) {
	if len(path) > len(t.root.path) {
		if path[:len(t.root.path)] != t.root.path {
			return zero
		}

		path = path[len(t.root.path):]

		return t.root.getFromChild(path, params)

	} else if path == t.root.path {
		switch {
		case t.root.handler != zero:
			return t.root.handler
		case t.root.wildcard != nil:
			return t.root.wildcard.handler
		}
	}

	return zero
}
