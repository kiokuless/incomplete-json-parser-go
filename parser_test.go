package incompletejson

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestIncompleteJsonParser_CompleteObjects(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"}`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_IncompleteObjects(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_IncompleteArrays(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `["apple","banana","orange"`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := []interface{}{"apple", "banana", "orange"}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_IncompleteStrings(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","message":"Hello, world!`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name":    "John",
		"message": "Hello, world!",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_MultipleChunks(t *testing.T) {
	parser := NewIncompleteJsonParser()
	chunk1 := `{"name":"John","a`
	chunk2 := `ge":30,"city":"New York"}`
	
	err := parser.Write(chunk1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	err = parser.Write(chunk2)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_NullValues(t *testing.T) {
	testCases := []string{"n", "nu", "nul", "null"}
	
	for _, nullValue := range testCases {
		parser := NewIncompleteJsonParser()
		jsonString := `{"name":"John","age":30,"isStudent":` + nullValue
		
		err := parser.Write(jsonString)
		if err != nil {
			t.Fatalf("Expected no error for %s, got %v", nullValue, err)
		}
		
		result, err := parser.GetObjects()
		if err != nil {
			t.Fatalf("Expected no error for %s, got %v", nullValue, err)
		}
		
		expected := map[string]interface{}{
			"name":      "John",
			"age":       float64(30),
			"isStudent": nil,
		}
		
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("For %s: Expected %+v, got %+v", nullValue, expected, result)
		}
	}
}

func TestIncompleteJsonParser_NestedObjects(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"address":{"street":"123 Main St","city":"New York"`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"address": map[string]interface{}{
			"street": "123 Main St",
			"city":   "New York",
		},
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_Parse(t *testing.T) {
	jsonString := `{"name":"John","age":30,"city":"New York"}`
	
	result, err := Parse(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_Reset(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New `
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New ",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
	
	parser.Reset()
	_, err = parser.GetObjects()
	if err == nil {
		t.Error("Expected error after reset, got nil")
	}
}

func TestIncompleteJsonParser_TrailingWhitespaces(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"}  
  `
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_RedundantComma(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York",`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_SameValueRepeated(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result1, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result2, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": "New York",
	}
	
	if !reflect.DeepEqual(result1, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result1)
	}
	
	if !reflect.DeepEqual(result2, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result2)
	}
}

func TestIncompleteJsonParser_EndWithColon(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": nil,
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_IncrementalKey(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"cit`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"cit":  nil,
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
	
	// Continue writing
	err = parser.Write("y")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err = parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected = map[string]interface{}{
		"name": "John",
		"age":  float64(30),
		"city": nil,
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
	
	// Add quote and colon
	err = parser.Write(`":`)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err = parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
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
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
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
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_ComplexNestedWithArray(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"address":{"street":"123 Main St","city":"New York","zip":10001, "alias": ["Dante"`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
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
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
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
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	// Expected should have actual newlines, tabs, etc.
	expected := map[string]interface{}{
		"message": "Hello\nWorld\tTab\r\nNewline\\Backslash\"Quote",
		"simple":  "No escapes here",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
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
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"text": "Hello\nHello",
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}

func TestIncompleteJsonParser_EscapedCharactersInStrings(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","message":"Hello, \"World\"! [{}]"}`
	
	err := parser.Write(jsonString)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	result, err := parser.GetObjects()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	
	expected := map[string]interface{}{
		"name":    "John",
		"message": `Hello, "World"! [{}]`,
	}
	
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %+v, got %+v", expected, result)
	}
}