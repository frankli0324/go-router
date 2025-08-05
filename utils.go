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

func lcp(a, b string) int {
	for i, ra := range a {
		if i >= len(b) {
			return i
		}
		if rb, _ := utf8.DecodeRuneInString(b[i:]); ra != rb {
			return i
		}
	}
	return len(a)
}
