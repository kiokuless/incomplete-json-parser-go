package incompletejson

import (
	"unicode"
)

// Scope interface defines the contract for all scope types
type Scope interface {
	Write(letter rune) bool
	GetOrAssume() interface{}
	IsFinished() bool
	SetAllowUnescapedNewlines(allow bool)
}

// BaseScope provides common functionality for all scopes
type BaseScope struct {
	finish                 bool
	allowUnescapedNewlines bool
}

func (s *BaseScope) IsFinished() bool {
	return s.finish
}

func (s *BaseScope) SetAllowUnescapedNewlines(allow bool) {
	s.allowUnescapedNewlines = allow
}

// isWhitespace checks if a character is whitespace
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
