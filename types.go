package router

import (
	"regexp"
	"strings"
)

type holder struct{ k, v string }
type matcher interface {
	match(string, *holder) (int, bool)
	equal(matcher) bool
	string() string
}

type literal string

func (l literal) match(s string, _ *holder) (i int, ok bool) {
	i = len(l)
	ok = len(s) >= i && s[:i] == string(l)
	if !ok {
		i = 0
	}
	return
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

func (p param) match(s string, h *holder) (int, bool) {
	if s == "" {
		return 0, false
	}
	i := strings.IndexByte(s, '/')
	if i == -1 {
		i = len(s)
	}
	if p.after != "" {
		i = strings.Index(s[:i], p.after)
	}
	if i == -1 {
		return 0, false
	}
	*h = holder{p.key, s[:i]}
	return i, true
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

func (w wildcard) match(s string, h *holder) (int, bool) {
	*h = holder{string(w), s}
	return len(s), true
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

func (w regex) match(s string, h *holder) (int, bool) {
	match := w.FindString(s)
	if match == "" {
		return 0, false
	}
	*h = holder{w.key, match}
	return len(match), true
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
