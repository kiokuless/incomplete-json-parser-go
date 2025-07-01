package incompletejson

import (
	"errors"
)

// IncompleteJsonParser is the main parser struct
type IncompleteJsonParser struct {
	scope  Scope
	finish bool
}

// NewIncompleteJsonParser creates a new parser instance
func NewIncompleteJsonParser() *IncompleteJsonParser {
	return &IncompleteJsonParser{}
}

// Parse is a static method that parses JSON in a single step
func Parse(chunk string) (interface{}, error) {
	parser := NewIncompleteJsonParser()
	err := parser.Write(chunk)
	if err != nil {
		return nil, err
	}
	return parser.GetObjects()
}

// Reset resets the parser's internal state
func (p *IncompleteJsonParser) Reset() {
	p.scope = nil
	p.finish = false
}

// Write processes a chunk of JSON data
func (p *IncompleteJsonParser) Write(chunk string) error {
	for _, letter := range chunk {
		if p.finish {
			if isWhitespace(letter) {
				continue
			}
			return errors.New("parser is already finished")
		}

		if p.scope == nil {
			if isWhitespace(letter) {
				continue
			} else if letter == '{' {
				p.scope = NewObjectScope()
			} else if letter == '[' {
				p.scope = NewArrayScope()
			} else {
				p.scope = NewLiteralScope()
			}
			success := p.scope.Write(letter)
			if !success {
				return errors.New("failed to parse the JSON string")
			}
		} else {
			success := p.scope.Write(letter)
			if success {
				if p.scope.IsFinished() {
					p.finish = true
					continue
				}
			} else {
				return errors.New("failed to parse the JSON string")
			}
		}
	}
	return nil
}

// GetObjects returns the parsed JavaScript object
func (p *IncompleteJsonParser) GetObjects() (interface{}, error) {
	if p.scope != nil {
		return p.scope.GetOrAssume(), nil
	}
	return nil, errors.New("no input to parse")
}