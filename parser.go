package incompletejson

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// IncompleteJsonParser is the main parser struct
type IncompleteJsonParser struct {
	scope                  Scope
	finish                 bool
	ignoreExtraCharacters  bool
	allowUnescapedNewlines bool
	validateRequiredFields bool
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

// WithRequiredFields sets the option to validate that all non-omitempty fields are present
func WithRequiredFields(validate bool) ParserOption {
	return func(p *IncompleteJsonParser) {
		p.validateRequiredFields = validate
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
func Parse(chunk string, options ...ParserOption) (interface{}, error) {
	parser := NewIncompleteJsonParser(options...)
	err := parser.Write(chunk)
	if err != nil {
		return nil, err
	}
	return parser.GetObjects()
}

// UnmarshalTo is a static method that parses JSON and stores the result in the value pointed to by v
func UnmarshalTo(chunk string, v interface{}, options ...ParserOption) error {
	parser := NewIncompleteJsonParser(options...)
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

	// If result is nil (null JSON), return an error for type safety
	if result == nil {
		return errors.New("cannot unmarshal null into struct")
	}

	// Convert to JSON bytes and then unmarshal to the target type
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, v)
	if err != nil {
		return err
	}

	// Validate required fields if option is enabled
	if p.validateRequiredFields {
		return p.validateRequired(v, result)
	}

	return nil
}

// GetObjectsAs returns the parsed data as the specified type using generics
func GetObjectsAs[T any](p *IncompleteJsonParser) (T, error) {
	var result T
	err := p.UnmarshalTo(&result)
	return result, err
}

// ParseAs is a static generic function that parses JSON and returns the result as the specified type
func ParseAs[T any](chunk string, options ...ParserOption) (T, error) {
	var result T
	err := UnmarshalTo(chunk, &result, options...)
	return result, err
}

// validateRequired checks that all non-omitempty fields are present in the JSON
func (p *IncompleteJsonParser) validateRequired(target interface{}, jsonData interface{}) error {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() == reflect.Ptr {
		targetValue = targetValue.Elem()
	}

	if targetValue.Kind() != reflect.Struct {
		// Not a struct, no validation needed
		return nil
	}

	targetType := targetValue.Type()
	jsonMap, ok := jsonData.(map[string]interface{})
	if !ok {
		return fmt.Errorf("expected JSON object for struct type %s", targetType.Name())
	}

	var missingFields []string

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		jsonTag := field.Tag.Get("json")

		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// Parse the tag to get field name and check for omitempty
		tagParts := strings.Split(jsonTag, ",")
		fieldName := tagParts[0]
		hasOmitEmpty := false

		for j := 1; j < len(tagParts); j++ {
			if tagParts[j] == "omitempty" {
				hasOmitEmpty = true
				break
			}
		}

		// If field doesn't have omitempty and is not present in JSON, it's an error
		if !hasOmitEmpty {
			if _, exists := jsonMap[fieldName]; !exists {
				missingFields = append(missingFields, fieldName)
			}
		}
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missingFields, ", "))
	}

	return nil
}
