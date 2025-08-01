package router

import (
	"regexp"
	"strings"
)

type matcher interface {
	match(string) (int, string, bool)
	equal(matcher) bool
	string() string
}

type literal string

func (l literal) match(s string) (int, string, bool) {
	return len(l), "", len(s) >= len(l) && s[:len(l)] == string(l)
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

func (p param) match(s string) (int, string, bool) {
	if s == "" {
		return 0, "", false
	}
	i := strings.IndexByte(s, '/')
	if i == -1 {
		i = len(s)
	}
	if p.after != "" {
		i = strings.Index(s[:i], p.after)
	}
	return i, p.key, i != -1
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

func (w wildcard) match(s string) (int, string, bool) {
	return len(s), string(w), true
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

func (w regex) match(s string) (int, string, bool) {
	match := w.FindString(s)
	return len(match), w.key, match != ""
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
