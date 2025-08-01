package radix

import (
	"sort"
	"strings"

	gstrings "github.com/savsgio/gotils/strings"
)

func newNode[T comparable](path string) *node[T] {
	return &node[T]{
		nType: static,
		path:  path,
	}
}

// conflict raises a panic with some details
func (n *nodeWildcard[T]) conflict(path, fullPath string) error {
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildcardConflict, path, fullPath, n.path, prefix)
}

// wildPathConflict raises a panic with some details
func (n *node[T]) wildPathConflict(path, fullPath string) error {
	pathSeg := strings.SplitN(path, "/", 2)[0]
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildPathConflict, pathSeg, fullPath, n.path, prefix)
}

// clone clones the current node in a new pointer
func (n node[T]) clone() *node[T] {
	cloneNode := new(node[T])
	cloneNode.nType = n.nType
	cloneNode.path = n.path
	cloneNode.handler = n.handler

	if len(n.children) > 0 {
		cloneNode.children = make([]*node[T], len(n.children))

		for i, child := range n.children {
			cloneNode.children[i] = child.clone()
		}
	}

	if n.wildcard != nil {
		cloneNode.wildcard = &nodeWildcard[T]{
			path:     n.wildcard.path,
			paramKey: n.wildcard.paramKey,
			handler:  n.wildcard.handler,
		}
	}

	if len(n.paramKeys) > 0 {
		cloneNode.paramKeys = make([]string, len(n.paramKeys))
		copy(cloneNode.paramKeys, n.paramKeys)
	}

	cloneNode.paramRegex = n.paramRegex

	return cloneNode
}

func (n *node[T]) split(i int) {
	cloneChild := n.clone()
	cloneChild.nType = static
	cloneChild.path = cloneChild.path[i:]
	cloneChild.paramKeys = nil
	cloneChild.paramRegex = nil

	var zero T
	n.path = n.path[:i]
	n.handler = zero
	n.wildcard = nil
	n.children = append(n.children[:0], cloneChild)
}

func (n *node[T]) findEndIndexAndValues(path string) (int, []string) {
	index := n.paramRegex.FindStringSubmatchIndex(path)
	if len(index) == 0 || index[0] != 0 {
		return -1, nil
	}

	end := index[1]

	index = index[2:]
	values := make([]string, len(index)/2)

	i := 0
	for j := range index {
		if (j+1)%2 != 0 {
			continue
		}

		values[i] = gstrings.Copy(path[index[j-1]:index[j]])

		i++
	}

	return end, values
}

func (n *node[T]) setHandler(handler T, fullPath string) (*node[T], error) {
	var zero T
	if n.handler != zero {
		return n, newRadixError(errSetHandler, fullPath)
	}

	n.handler = handler
	return n, nil
}

func (n *node[T]) insert(path, fullPath string, handler T) (*node[T], error) {
	end := segmentEndIndex(path, true)
	child := newNode[T](path)

	wp := findWildPath(path, fullPath)
	if wp != nil {
		j := end
		if wp.start > 0 {
			j = wp.start
		}

		child.path = path[:j]

		if wp.start > 0 {
			n.children = append(n.children, child)

			return child.insert(path[j:], fullPath, handler)
		}

		switch wp.pType {
		case param:
			n.hasWildChild = true

			child.nType = wp.pType
			child.paramKeys = wp.keys
			child.paramRegex = wp.regex
		case wildcard:
			if len(path) == end && n.path[len(n.path)-1] != '/' {
				return nil, newRadixError(errWildcardSlash, fullPath)
			} else if len(path) != end {
				return nil, newRadixError(errWildcardNotAtEnd, fullPath)
			}

			if n.wildcard != nil {
				if n.wildcard.path == path {
					return n, newRadixError(errSetWildcardHandler, fullPath)
				}

				return nil, n.wildcard.conflict(path, fullPath)
			}

			n.wildcard = &nodeWildcard[T]{
				path:     wp.path,
				paramKey: wp.keys[0],
				handler:  handler,
			}

			return n, nil
		}

		path = path[wp.end:]

		if len(path) > 0 {
			n.children = append(n.children, child)

			return child.insert(path, fullPath, handler)
		}
	}

	child.handler = handler
	n.children = append(n.children, child)

	return child, nil
}

// add adds the handler to node for the given path
func (n *node[T]) add(path, fullPath string, handler T) (*node[T], error) {
	var zero T
	if len(path) == 0 {
		return n.setHandler(handler, fullPath)
	}

	for _, child := range n.children {
		i := longestCommonPrefix(path, child.path)
		if i == 0 {
			continue
		}

		switch child.nType {
		case static:
			if len(child.path) > i {
				child.split(i)
			}

			if len(path) > i {
				return child.add(path[i:], fullPath, handler)
			}
		case param:
			wp := findWildPath(path, fullPath)

			isParam := wp.start == 0 && wp.pType == param
			hasHandler := child.handler != zero || handler == zero

			if len(path) == wp.end && isParam && hasHandler {
				// The current segment is a param and it's duplicated
				if child.path == path {
					return child, newRadixError(errSetHandler, fullPath)
				}

				return nil, child.wildPathConflict(path, fullPath)
			}

			if len(path) > i {
				if child.path == wp.path {
					return child.add(path[i:], fullPath, handler)
				}

				return n.insert(path, fullPath, handler)
			}
		}

		return child.setHandler(handler, fullPath)
	}

	return n.insert(path, fullPath, handler)
}

func (n *node[T]) getFromChild(path string, params map[string]string) (zero T) {
	for _, child := range n.children {
		switch child.nType {
		case static:

			// Checks if the first byte is equal
			// It's faster than compare strings
			if path[0] != child.path[0] {
				continue
			}

			if len(path) > len(child.path) {
				if path[:len(child.path)] != child.path {
					continue
				}

				if h := child.getFromChild(path[len(child.path):], params); h != zero {
					return h
				}
			} else if path == child.path {
				switch {
				case child.handler != zero:
					return child.handler
				case child.wildcard != nil:
					params[child.wildcard.paramKey] = ""

					return child.wildcard.handler
				}

				return zero
			}

		case param:
			end := segmentEndIndex(path, true)
			values := []string{gstrings.Copy(path[:end])}

			if child.paramRegex != nil {
				end, values = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			if len(path) > end {
				if h := child.getFromChild(path[end:], params); h != zero {
					for i, key := range child.paramKeys {
						params[key] = values[i]
					}

					return h
				}

			} else if len(path) == end {
				switch {
				case child.handler == zero:
					// try another child
					continue
				}
				for i, key := range child.paramKeys {
					params[key] = values[i]
				}

				return child.handler
			}

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		params[n.wildcard.paramKey] = gstrings.Copy(path)
		return n.wildcard.handler
	}

	return zero
}

// sort sorts the current node and their children
func (n *node[T]) sort() {
	for _, child := range n.children {
		child.sort()
	}

	sort.Sort(n)
}

// Len returns the total number of children the node has
func (n *node[T]) Len() int {
	return len(n.children)
}

// Swap swaps the order of children nodes
func (n *node[T]) Swap(i, j int) {
	n.children[i], n.children[j] = n.children[j], n.children[i]
}

// Less checks if the node 'i' has less priority than the node 'j'
func (n *node[T]) Less(i, j int) bool {
	if n.children[i].nType < n.children[j].nType {
		return true
	} else if n.children[i].nType > n.children[j].nType {
		return false
	}

	return len(n.children[i].children) > len(n.children[j].children)
}
