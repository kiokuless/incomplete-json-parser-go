package incompletejson

// ArrayScope handles parsing of JSON arrays
type ArrayScope struct {
	BaseScope
	array []Scope
	state string // "value" or "comma"
	scope Scope
}

func NewArrayScope() *ArrayScope {
	return &ArrayScope{
		array: make([]Scope, 0),
		state: "value",
	}
}

func (a *ArrayScope) Write(letter rune) bool {
	if a.finish {
		return false
	}

	// Ignore first [
	if len(a.array) == 0 && a.state == "value" && a.scope == nil {
		if letter == '[' {
			return true
		}
	}

	// Process the letter
	switch a.state {
	case "value":
		if a.scope == nil {
			if isWhitespace(letter) {
				return true
			} else if letter == ']' {
				// Empty array case: []
				a.finish = true
				return true
			} else if letter == '{' {
				a.scope = NewObjectScope()
				a.scope.SetAllowUnescapedNewlines(a.allowUnescapedNewlines)
				a.array = append(a.array, a.scope)
				return a.scope.Write(letter)
			} else if letter == '[' {
				a.scope = NewArrayScope()
				a.scope.SetAllowUnescapedNewlines(a.allowUnescapedNewlines)
				a.array = append(a.array, a.scope)
				return a.scope.Write(letter)
			} else {
				a.scope = NewLiteralScope()
				a.scope.SetAllowUnescapedNewlines(a.allowUnescapedNewlines)
				a.array = append(a.array, a.scope)
				return a.scope.Write(letter)
			}
		} else {
			success := a.scope.Write(letter)
			if success {
				if a.scope.IsFinished() {
					a.state = "comma"
				}
				return true
			} else {
				if a.scope.IsFinished() {
					a.state = "comma"
					return true
				} else if letter == ',' {
					a.scope = nil
				} else if letter == ']' {
					a.finish = true
					return true
				}
				return true
			}
		}
	case "comma":
		if isWhitespace(letter) {
			return true
		} else if letter == ',' {
			a.state = "value"
			a.scope = nil
			return true
		} else if letter == ']' {
			a.finish = true
			return true
		} else {
			return false
		}
	}

	return false
}

func (a *ArrayScope) GetOrAssume() interface{} {
	result := make([]interface{}, len(a.array))
	for i, scope := range a.array {
		result[i] = scope.GetOrAssume()
	}
	return result
}
