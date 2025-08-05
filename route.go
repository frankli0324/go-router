package router

type Router[T any] struct {
	tree node[T]
}

// Set registers a value for the given URL pattern. It's not routine-safe.
func (r *Router[T]) Set(path string, handler T) error {
	if path == "" || path[0] != '/' {
		return ErrInvalidPath.With(path)
	}
	_, err := r.tree.add(path, path, handler)
	r.tree.sort()
	return err
}

// GetParam matches the given path and returns the corresponding value,
// assigning the given params map with the matched parameters.
// If no pattern is found, the zero value is returned. It's routine-safe.
func (r *Router[T]) GetParam(path string, params map[string]string) (zero T) {
	if path == "" || path[0] != '/' || len(r.tree.children) == 0 {
		return zero
	}
	n := &r.tree
	if ch := n.children[0]; ch.b == '/' {
		// first node is almost always a literal("/")
		m := ch.m.(literal)
		end, ok := 1, len(m) == 1
		if !ok {
			end, _, ok = m.match(path)
		}
		if ok {
			n, path = ch, path[end:]
			if path == "" {
				zero = n.handler
			}
		}
	}
	n = n.get(path, params)
	if n == nil {
		return zero
	}
	return n.handler
}

// Get matches the given path and returns the corresponding value.
// If no pattern is found, the zero value is returned. It's routine-safe.
func (r *Router[T]) Get(path string) T {
	return r.GetParam(path, nil)
}

func NewRouter[T any]() *Router[T] {
	return &Router[T]{tree: node[T]{m: literal("")}}
}
