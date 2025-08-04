package router

import (
	"sort"
	"strings"
)

func (n *node[T]) add(path, fullPath string, handler T) (*node[T], error) {
	for path != "" {
		next, end, err := next(path)
		if err != nil {
			return nil, err
		}
		inserted := false
		for _, child := range n.children {
			if child.m.equal(next) {
				n = child
				inserted = true
				path = path[end:]
				break
			}
		}

		switch next := next.(type) {
		case literal:
			maxi, maxl := -1, 0
			for i, child := range n.children {
				l, ok := child.m.(literal)
				if !ok || string(l)[0] != path[0] {
					continue
				}
				if common := lcp(string(l), string(next)); common > maxl {
					maxi, maxl = i, common
				}
			}
			if maxi != -1 {
				// split node
				child := n.children[maxi]
				n = child.cut(maxl)
				path = path[maxl:]
				inserted = true
				break
			}
		case wildcard, param:
			for _, child := range n.children {
				if typeID(child.m) == typeID(next) {
					return nil, ErrExprConflict.With(fullPath, child.m.string())
				}
			}
		}
		if inserted {
			continue
		}
		path = path[end:]
		if _, ok := next.(wildcard); ok && path != "" {
			return nil, ErrWildcardNotAtEnd.With(fullPath)
		}
		newch := &node[T]{m: next}
		if next, ok := next.(literal); ok {
			newch.b = next[0]
		}
		n.children = append(n.children, newch)
		n = newch
	}
	if n.assigned {
		return nil, ErrConflict.With(fullPath)
	}
	n.handler = handler
	n.assigned = true
	return n, nil
}

func (n *node[T]) get(path string, params map[string]string) *node[T] {
	for _, child := range n.children {
		end, ok, key := 0, false, ""
		if child.b != 0 {
			if path != "" && child.b != path[0] {
				continue
			}
			end, key, ok = child.m.(literal).match(path)
		} else {
			end, key, ok = child.m.match(path)
		}
		if !ok {
			continue
		}
		next := child.get(path[end:], params)
		if next == nil {
			if end != len(path) {
				continue
			}
			next = child
		}
		if params != nil && next != nil && next.assigned && key != "" {
			params[key] = strings.Clone(path[:end])
		}
		return next
	}
	return nil
}

func (n *node[T]) cut(i int) *node[T] {
	l, ok := n.m.(literal)
	if !ok {
		panic("cannot cut non-literal node")
	}
	if i == len(l) {
		return n
	}
	n.children = []*node[T]{{
		m: l[i:], b: l[i], children: n.children,
		handler: n.handler, assigned: n.assigned,
	}}
	var zero T
	n.handler = zero
	n.assigned = false
	n.m = l[:i]
	return n
}

func (n *node[T]) sort() {
	sort.Sort(n)
	for _, child := range n.children {
		child.sort()
	}
}
