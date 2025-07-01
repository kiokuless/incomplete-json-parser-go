package incompletejson

import (
	"unicode"
)

// Scope interface defines the contract for all scope types
type Scope interface {
	Write(letter rune) bool
	GetOrAssume() interface{}
	IsFinished() bool
}

// BaseScope provides common functionality for all scopes
type BaseScope struct {
	finish bool
}

func (s *BaseScope) IsFinished() bool {
	return s.finish
}

// isWhitespace checks if a character is whitespace
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}
