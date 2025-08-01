package router

import "fmt"

type err struct {
	string
	params []interface{}
}

func (e *err) Error() string {
	return fmt.Sprintf(e.string, e.params...)
}

func (e *err) With(params ...interface{}) *err {
	return &err{e.string, params}
}

var (
	ErrInvalidPath      = &err{"path must begin with '/' in path '%s'", nil}
	ErrExprConflict     = &err{"path '%s' conflicts with existing wildcard or param '%s'", nil}
	ErrConflict         = &err{"a handler is already registered for path '%s'", nil}
	ErrExpr             = &err{"invalid expression '%s': '%s'", nil}
	ErrWildcardNotAtEnd = &err{"wildcard routes are only allowed at the end of the path in path '%s'", nil}
)
