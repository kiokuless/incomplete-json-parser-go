package incompletejson

// ObjectScope handles parsing of JSON objects
type ObjectScope struct {
	BaseScope
	object     map[string]interface{}
	state      string // "key", "colons", "value", "comma"
	keyScope   *LiteralScope
	valueScope Scope
}

func NewObjectScope() *ObjectScope {
	return &ObjectScope{
		object: make(map[string]interface{}),
		state:  "key",
	}
}

func (o *ObjectScope) Write(letter rune) bool {
	if o.finish {
		return false
	}

	// Ignore first {
	if len(o.object) == 0 && o.state == "key" && o.keyScope == nil && o.valueScope == nil {
		if letter == '{' {
			return true
		}
	}

	switch o.state {
	case "key":
		if o.keyScope == nil {
			if isWhitespace(letter) {
				return true
			} else if letter == '}' {
				// Empty object case: {}
				o.finish = true
				return true
			} else if letter == '"' {
				o.keyScope = NewLiteralScope()
				o.keyScope.SetAllowUnescapedNewlines(o.allowUnescapedNewlines)
				return o.keyScope.Write(letter)
			} else {
				return false
			}
		} else {
			o.keyScope.Write(letter)
			key := o.keyScope.GetOrAssume()
			if _, ok := key.(string); ok {
				if o.keyScope.IsFinished() {
					o.state = "colons"
				}
				return true
			} else {
				return false
			}
		}

	case "colons":
		if isWhitespace(letter) {
			return true
		} else if letter == ':' {
			o.state = "value"
			o.valueScope = nil
			return true
		} else {
			return false
		}

	case "value":
		if o.valueScope == nil {
			if isWhitespace(letter) {
				return true
			} else if letter == '{' {
				o.valueScope = NewObjectScope()
				o.valueScope.SetAllowUnescapedNewlines(o.allowUnescapedNewlines)
				return o.valueScope.Write(letter)
			} else if letter == '[' {
				o.valueScope = NewArrayScope()
				o.valueScope.SetAllowUnescapedNewlines(o.allowUnescapedNewlines)
				return o.valueScope.Write(letter)
			} else {
				o.valueScope = NewLiteralScope()
				o.valueScope.SetAllowUnescapedNewlines(o.allowUnescapedNewlines)
				return o.valueScope.Write(letter)
			}
		} else {
			success := o.valueScope.Write(letter)
			if o.valueScope.IsFinished() {
				key := o.keyScope.GetOrAssume().(string)
				o.object[key] = o.valueScope.GetOrAssume()
				o.state = "comma"
				return true
			} else if success {
				return true
			} else {
				if isWhitespace(letter) {
					return true
				} else if letter == ',' {
					key := o.keyScope.GetOrAssume().(string)
					o.object[key] = o.valueScope.GetOrAssume()
					o.state = "key"
					o.keyScope = nil
					o.valueScope = nil
					return true
				} else if letter == '}' {
					key := o.keyScope.GetOrAssume().(string)
					o.object[key] = o.valueScope.GetOrAssume()
					o.finish = true
					return true
				} else {
					return false
				}
			}
		}

	case "comma":
		if isWhitespace(letter) {
			return true
		} else if letter == ',' {
			o.state = "key"
			o.keyScope = nil
			o.valueScope = nil
			return true
		} else if letter == '}' {
			o.finish = true
			return true
		} else {
			return false
		}
	}

	return false
}

func (o *ObjectScope) GetOrAssume() interface{} {
	result := make(map[string]interface{})

	// Copy existing completed key-value pairs
	for k, v := range o.object {
		result[k] = v
	}

	// Handle incomplete key-value pair
	if o.keyScope != nil || o.valueScope != nil {
		if o.keyScope != nil {
			key := o.keyScope.GetOrAssume()
			if keyStr, ok := key.(string); ok && len(keyStr) > 0 {
				var value interface{}
				if o.valueScope != nil {
					value = o.valueScope.GetOrAssume()
				}
				if value != nil {
					result[keyStr] = value
				} else {
					result[keyStr] = nil
				}
			}
		}
	}

	return result
}
