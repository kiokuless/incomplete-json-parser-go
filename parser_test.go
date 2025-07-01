package incompletejson

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIncompleteJsonParser_CompleteObjects(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"}`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_IncompleteObjects(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_IncompleteArrays(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `["apple","banana","orange"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := []interface{}{"apple", "banana", "orange"}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_IncompleteStrings(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","message":"Hello, world!`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name":    "John",
		"message": "Hello, world!",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_MultipleChunks(t *testing.T) {
	parser := NewIncompleteJsonParser()
	chunk1 := `{"name":"John","a`
	chunk2 := `ge":30,"city":"New York"}`

	err := parser.Write(chunk1)
	require.NoError(t, err)

	err = parser.Write(chunk2)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_NullValues(t *testing.T) {
	testCases := []string{"n", "nu", "nul", "null"}

	for _, nullValue := range testCases {
		t.Run(nullValue, func(t *testing.T) {
			parser := NewIncompleteJsonParser()
			jsonString := `{"name":"John","age":30,"isStudent":` + nullValue

			err := parser.Write(jsonString)
			require.NoError(t, err)

			result, err := parser.GetObjects()
			require.NoError(t, err)

			expected := map[string]interface{}{
				"name":      "John",
				"age":       float64(30),
				"isStudent": nil,
			}

			require.Equal(t, expected, result)
		})
	}
}

func TestIncompleteJsonParser_NestedObjects(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"address":{"street":"123 Main St","city":"New York"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
		},
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_Parse(t *testing.T) {
	jsonString := `{"name":"John","age":30,"city":"New York"}`

	result, err := Parse(jsonString)
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_Reset(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New `

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New ",
	}

	require.Equal(t, expected, result)

	parser.Reset()
	_, err = parser.GetObjects()
	require.Error(t, err)
}

func TestIncompleteJsonParser_TrailingWhitespaces(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"}  
  `

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_RedundantComma(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York",`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_SameValueRepeated(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result1, err := parser.GetObjects()
	require.NoError(t, err)

	result2, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}

	require.Equal(t, expected, result1)
	require.Equal(t, expected, result2)
}

func TestIncompleteJsonParser_EndWithColon(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": nil,
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_IncrementalKey(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"cit`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"cit":  nil,
	}

	require.Equal(t, expected, result)

	// Continue writing
	err = parser.Write("y")
	require.NoError(t, err)

	result, err = parser.GetObjects()
	require.NoError(t, err)

	expected = map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": nil,
	}

	require.Equal(t, expected, result)

	// Add quote and colon
	err = parser.Write(`":`)
	require.NoError(t, err)

	result, err = parser.GetObjects()
	require.NoError(t, err)

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_ComplexNestedObject(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{
  "id": 12345,
  "name": "John Doe",
  "isActive": true,
  "age": 30,
  "email": "john.doe@example.com",
  "address": {
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zipCode": "10001",
    "country": "USA",
    "location": {
      "lat": 40.7128,
      "lng": -74.0060
    }
  }`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"id":       float64(12345),
		"name":     "John Doe",
		"isActive": true,
		"age":      float64(30),
		"email":    "john.doe@example.com",
		"address": map[string]interface{}{
			"street":  "123 Main St",
			"city":    "New York",
			"state":   "NY",
			"zipCode": "10001",
			"country": "USA",
			"location": map[string]interface{}{
				"lat": 40.7128,
				"lng": -74.006,
			},
		},
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_ComplexNestedWithArray(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"address":{"street":"123 Main St","city":"New York","zip":10001, "alias": ["Dante"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
			"zip":    float64(10001),
			"alias":  []interface{}{"Dante"},
		},
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_EscapeCharacters(t *testing.T) {
	parser := NewIncompleteJsonParser()

	// Create object with escape characters - matching TypeScript test
	obj := map[string]interface{}{
		"message": "Hello\nWorld\tTab\r\nNewline\\Backslash\"Quote",
		"simple":  "No escapes here",
	}

	// Convert to JSON string to get proper escaping
	jsonBytes, _ := json.Marshal(obj)
	jsonString := string(jsonBytes)

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	// Expected should have actual newlines, tabs, etc.
	expected := map[string]interface{}{
		"message": "Hello\nWorld\tTab\r\nNewline\\Backslash\"Quote",
		"simple":  "No escapes here",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_UnicodeCharacters(t *testing.T) {
	parser := NewIncompleteJsonParser()

	// Create object with Unicode characters
	obj := map[string]interface{}{
		"text": "\u0048\u0065\u006C\u006C\u006F\n\u0048\u0065\u006C\u006C\u006F",
	}

	// Convert to JSON string to get proper escaping
	jsonBytes, _ := json.Marshal(obj)
	jsonString := string(jsonBytes)

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"text": "Hello\nHello",
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_EscapedCharactersInStrings(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","message":"Hello, \"World\"! [{}]"}`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"name":    "John",
		"message": `Hello, "World"! [{}]`,
	}

	require.Equal(t, expected, result)
}

func TestIncompleteJsonParser_IgnoreExtraCharacters(t *testing.T) {
	parser := NewIncompleteJsonParser(WithIgnoreExtraCharacters(true))
	input := `{"response_text": "答え"}\n\nこれは追加の説明です。`

	err := parser.Write(input)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"response_text": "答え",
	}

	require.Equal(t, expected, result)
}
