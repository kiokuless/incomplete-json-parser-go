package incompletejson

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// LiteralScope handles parsing of literal values (strings, numbers, booleans, null)
type LiteralScope struct {
	BaseScope
	content string
}

func NewLiteralScope() *LiteralScope {
	return &LiteralScope{}
}

func (l *LiteralScope) Write(letter rune) bool {
	if l.finish {
		return false
	}

	l.content += string(letter)

	// For null values, we need to check the actual return value differently
	// since assume can be nil for valid null values
	if l.content != "" && !l.canParseContent() {
		l.content = l.content[:len(l.content)-1]
		return false
	}

	// Check if literal is complete
	if (l.isCompletedString()) ||
		(l.content == "true" || l.content == "false") ||
		(l.content == "null") {
		l.finish = true
	}

	return true
}

func (l *LiteralScope) canParseContent() bool {
	// Empty content is valid (assumes null)
	if l.content == "" {
		return true
	}

	// Check if it's a valid partial/complete null
	if strings.HasPrefix("null", l.content) {
		return true
	}

	// Check if it's a valid partial/complete boolean
	if strings.HasPrefix("true", l.content) || strings.HasPrefix("false", l.content) {
		return true
	}

	// Check if it's a valid string (starts with quote)
	if strings.HasPrefix(l.content, `"`) {
		return true
	}

	// Check if it's a valid number
	if l.content == "-" {
		return true
	}

	numberRegex := regexp.MustCompile(`^-?\d+(\.\d*)?$`)
	return numberRegex.MatchString(l.content)
}

func (l *LiteralScope) isCompletedString() bool {
	if !strings.HasPrefix(l.content, `"`) {
		return false
	}

	if len(l.content) < 2 {
		return false
	}

	// Check if string ends with unescaped quote
	// Count consecutive backslashes before the final quote
	if strings.HasSuffix(l.content, `"`) {
		// Check how many backslashes precede the quote
		escaped := false
		for i := len(l.content) - 2; i >= 0; i-- {
			if l.content[i] == '\\' {
				escaped = !escaped
			} else {
				break
			}
		}
		// If there's an odd number of backslashes, the quote is escaped
		return !escaped
	}

	return false
}

func (l *LiteralScope) GetOrAssume() interface{} {
	// Empty content assumes null
	if l.content == "" {
		return nil
	}

	// Null values
	if strings.HasPrefix("null", l.content) {
		return nil
	}

	// Boolean values
	if strings.HasPrefix("true", l.content) {
		return true
	}
	if strings.HasPrefix("false", l.content) {
		return false
	}

	// String values
	if strings.HasPrefix(l.content, `"`) {
		return l.parseString()
	}

	// Number values
	if l.content == "-" {
		return 0
	}

	numberRegex := regexp.MustCompile(`^-?\d+(\.\d*)?$`)
	if numberRegex.MatchString(l.content) {
		if num, err := strconv.ParseFloat(l.content, 64); err == nil {
			return num
		}
	}

	// Cannot assume
	return nil
}

func (l *LiteralScope) parseString() interface{} {
	jsonedString := l.content

	isCompleted := l.isCompletedString()

	if !isCompleted {
		// Remove incomplete unicode escape at the end
		unicodeRegex := regexp.MustCompile(`\\u[\da-fA-F]{0,3}$`)
		if unicodeRegex.MatchString(jsonedString) {
			match := unicodeRegex.FindStringIndex(jsonedString)
			if match != nil {
				jsonedString = jsonedString[:match[0]]
			}
		}

		// Remove meaningless backslash at the end (but not escaped backslash)
		for strings.HasSuffix(jsonedString, `\`) && !strings.HasSuffix(jsonedString, `\\`) {
			jsonedString = jsonedString[:len(jsonedString)-1]
		}

		jsonedString += `"`
	}

	var result string
	if err := json.Unmarshal([]byte(jsonedString), &result); err != nil {
		// Try to fix the string by removing problematic characters at the end
		for len(jsonedString) > 2 {
			if jsonedString[len(jsonedString)-2] == '\\' {
				// Remove the character before the closing quote and try again
				jsonedString = jsonedString[:len(jsonedString)-2] + `"`
				if err := json.Unmarshal([]byte(jsonedString), &result); err == nil {
					return result
				}
			} else {
				break
			}
		}
		// Silently return nil for unparseable strings
		return nil
	}

	return result
}
