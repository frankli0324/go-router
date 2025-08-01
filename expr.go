package router

import (
	"regexp"
	"strings"
)

// seg returns the index where the segment ends from the given path
func seg(path string) string {
	end := 0
	for end < len(path) && path[end] != '/' {
		end++
	}
	return path[:end]
}

func nextNonLiteral(path string) (matcher, int, error) {
	extend := -1
	keys := 0

	// Find end and check for invalid characters
	path = path[1:]
	for i, c := range []byte(path) {
		switch c {
		case '}':
			if len(path) > i+1 && path[i+1] == '{' {
				return nil, 0, ErrExpr.With("{"+path, "the expressions must be separated by at least 1 char")
			}
			ext := ""
			if extend == -1 {
				extend = i
			} else {
				ext = path[extend+1 : i]
				if ext == "" {
					return nil, 0, ErrExpr.With("{"+path, "empty match expression not allowed")
				}
			}
			key := path[:extend]
			if key == "" {
				return nil, 0, ErrExpr.With("{"+path, "wildcards must be named with a non-empty name")
			}
			switch ext {
			case "":
				s := seg(path[i+1:])
				return param{key: key, after: s}, len(s) + i + 2, nil
			case "*":
				return wildcard(key), i + 2, nil
			default:
				return regex{key, regexp.MustCompile(ext)}, i + 2, nil
			}
		case ':':
			extend = i
		case '{':
			if extend == -1 && keys == 0 {
				return nil, 0, ErrExpr.With(path, "the char '{' is not allowed in the param name")
			}

			keys++
		}
	}
	return nil, 0, nil
}

func next(path string) (matcher, int, error) {
	i := strings.IndexByte(path, '{')
	if i == -1 {
		return literal(path), len(path), nil
	}
	if i != 0 {
		return literal(path[:i]), i, nil
	}
	return nextNonLiteral(path)
}
