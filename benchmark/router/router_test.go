package router

import (
	"testing"

	"gotest.tools/v3/assert"
)

func v[T any](v T, _ error) T {
	return v
}

func TestRouter(t *testing.T) {
	t.Run("subpath matches", func(t *testing.T) {
		router := New[string]()
		router.Handle("/{path:*}", "/{path:*}")
		router.Handle("/test/{path:*}", "/test/{path:*}")
		assert.Equal(t, v(router.Handler("/asdfqwer")), "/{path:*}")
		assert.Equal(t, v(router.Handler("/test/asdfqwer")), "/test/{path:*}")
		assert.Equal(t, v(router.Handler("/test/asdf/qwerzxcv")), "/test/{path:*}")
	})
	t.Run("ignore trailing slash", func(t *testing.T) {
		router := New[string]()
		router.Handle("/{path:*}", "/{path:*}")
		router.Handle("/aabbcc/dde", "/aabbcc/dde")
		assert.Equal(t, v(router.Handler("/aabbcc/dde/")), "/aabbcc/dde")
		assert.Equal(t, v(router.Handler("/aabbcc/dde")), "/aabbcc/dde")
		assert.Equal(t, v(router.Handler("/aabbcc/def")), "/{path:*}")
	})
}
