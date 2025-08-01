package router

type Router[T any] struct {
	tree node[T]
}

func (r *Router[T]) Handle(path string, handler T) error {
	if path == "" || path[0] != '/' {
		return ErrInvalidPath.With(path)
	}
	_, err := r.tree.add(path, path, handler)
	r.tree.sort()
	return err
}

func (r *Router[T]) GetParam(path string, params map[string]string) (zero T) {
	if path == "" || path[0] != '/' {
		return zero
	}
	n := r.tree.get(path, params)
	if n == nil {
		return zero
	}
	return n.handler
}

func (r *Router[T]) Get(path string) T {
	return r.GetParam(path, nil)
}

func NewRouter[T any]() *Router[T] {
	return &Router[T]{tree: node[T]{m: literal("")}}
}
