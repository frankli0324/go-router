package router

import "unicode/utf8"

func typeID(m matcher) int {
	switch m.(type) {
	case literal:
		return 0
	case param:
		return 1
	case regex:
		return 2
	case wildcard:
		return 3
	}
	return -1
}

// Originates from fasthttp, licensed under the BSD 3-Clause License.
// https://github.com/fasthttp/router
func lcp(a, b string) int {
	i := 0

	max := func(a, b int) int {
		if a <= b {
			return a
		}
		return b
	}(utf8.RuneCountInString(a), utf8.RuneCountInString(b))

	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)
		rb, sizeB := utf8.DecodeRuneInString(b)

		a = a[sizeA:]
		b = b[sizeB:]

		if ra != rb {
			return i
		}

		i += sizeA
	}

	return i
}
