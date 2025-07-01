package incompletejson

import (
	"encoding/json"
	"errors"
)

// IncompleteJsonParser is the main parser struct
type IncompleteJsonParser struct {
	scope                  Scope
	finish                 bool
	ignoreExtraCharacters  bool
	allowUnescapedNewlines bool
}

// ParserOption defines a function type for parser options
type ParserOption func(*IncompleteJsonParser)

// WithIgnoreExtraCharacters sets the option to ignore extra characters after JSON completion
func WithIgnoreExtraCharacters(ignore bool) ParserOption {
	return func(p *IncompleteJsonParser) {
		p.ignoreExtraCharacters = ignore
	}
}

// WithAllowUnescapedNewlines sets the option to allow unescaped newlines in JSON strings
func WithAllowUnescapedNewlines(allow bool) ParserOption {
	return func(p *IncompleteJsonParser) {
		p.allowUnescapedNewlines = allow
	}
}

// NewIncompleteJsonParser creates a new parser instance with optional configuration
func NewIncompleteJsonParser(options ...ParserOption) *IncompleteJsonParser {
	parser := &IncompleteJsonParser{}

	// Apply all options
	for _, option := range options {
		option(parser)
	}

	return parser
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

// UnmarshalTo is a static method that parses JSON and stores the result in the value pointed to by v
func UnmarshalTo(chunk string, v interface{}) error {
	parser := NewIncompleteJsonParser()
	err := parser.Write(chunk)
	if err != nil {
		return err
	}
	return parser.UnmarshalTo(v)
}

// Reset resets the parser's internal state
func (p *IncompleteJsonParser) Reset() {
	p.scope = nil
	p.finish = false
	// ignoreExtraCharacters設定は保持する
}

// Write processes a chunk of JSON data
func (p *IncompleteJsonParser) Write(chunk string) error {
	for _, letter := range chunk {
		if p.finish {
			if p.ignoreExtraCharacters {
				// オプションが有効な場合は余分な文字を無視
				continue
			}
			// デフォルトの動作：空白文字のみ許可
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
			p.scope.SetAllowUnescapedNewlines(p.allowUnescapedNewlines)
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

// UnmarshalTo parses the JSON data and stores the result in the value pointed to by v
func (p *IncompleteJsonParser) UnmarshalTo(v interface{}) error {
	result, err := p.GetObjects()
	if err != nil {
		return err
	}

	// Convert to JSON bytes and then unmarshal to the target type
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonBytes, v)
}

// GetObjectsAs returns the parsed data as the specified type using generics
func GetObjectsAs[T any](p *IncompleteJsonParser) (T, error) {
	var result T
	err := p.UnmarshalTo(&result)
	return result, err
}

// ParseAs is a static generic function that parses JSON and returns the result as the specified type
func ParseAs[T any](chunk string) (T, error) {
	var result T
	err := UnmarshalTo(chunk, &result)
	return result, err
}
