package router

import (
	"regexp"
	"strings"
)

type matcher interface {
	match(string) (int, bool, func(v map[string]string))
	equal(matcher) bool
	string() string
}

type literal string

func (l literal) match(s string) (int, bool, func(v map[string]string)) {
	if len(l) > len(s) {
		return 0, false, nil
	}
	if s[:len(l)] == string(l) {
		return len(l), true, nil
	}
	return 0, false, nil
}

func (l literal) equal(m matcher) bool {
	if m2, ok := m.(literal); ok {
		return l == m2
	}
	return false
}

func (l literal) string() string {
	return string(l)
}

type param struct {
	key   string
	after string
}

func (p param) match(s string) (int, bool, func(v map[string]string)) {
	if s == "" {
		return 0, false, nil
	}
	var i int
	if p.after == "" {
		i = strings.IndexByte(s, '/')
		if i == -1 {
			i = len(s)
		}
	} else {
		i = strings.Index(s, p.after)
	}
	if i == -1 {
		return 0, false, nil
	}
	return i, true, func(v map[string]string) {
		v[p.key] = s[:i]
	}
}

func (p param) equal(m matcher) bool {
	if m2, ok := m.(param); ok {
		return p.key == m2.key && p.after == m2.after
	}
	return false
}

func (p param) string() string {
	return "{" + p.key + "}" + p.after
}

type wildcard string

func (w wildcard) match(s string) (int, bool, func(v map[string]string)) {
	return len(s), true, func(v map[string]string) {
		v[string(w)] = s
	}
}

func (w wildcard) equal(m matcher) bool {
	if m2, ok := m.(wildcard); ok {
		return w == m2
	}
	return false
}

func (w wildcard) string() string {
	return "{" + string(w) + ":*}"
}

type regex struct {
	key string
	*regexp.Regexp
}

func (w regex) match(s string) (int, bool, func(v map[string]string)) {
	matches := w.FindStringSubmatch(s)
	if matches == nil {
		return 0, false, nil
	}
	return len(matches[0]), true, func(v map[string]string) {
		v[w.key] = matches[0]
		for i, name := range w.SubexpNames()[1:] {
			v[name] = matches[i+1]
		}
	}
}

func (w regex) equal(m matcher) bool {
	if m2, ok := m.(regex); ok {
		return w.String() == m2.String()
	}
	return false
}

func (w regex) string() string {
	return "{" + w.key + ":" + w.Regexp.String() + "}"
}

type node[T any] struct {
	m        matcher
	assigned bool
	handler  T
	children []*node[T]
}
