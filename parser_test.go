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
	input := `{"message": "こんにちは世界"}\n\nこれは追加のテキストです。`

	err := parser.Write(input)
	require.NoError(t, err)

	result, err := parser.GetObjects()
	require.NoError(t, err)

	expected := map[string]interface{}{
		"message": "こんにちは世界",
	}

	require.Equal(t, expected, result)
}

// Test structures for type mapping
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	ZipCode string `json:"zipCode"`
}

type PersonWithAddress struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Address Address `json:"address"`
}

type BlogPost struct {
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
	Comments []string `json:"comments"`
}

func TestIncompleteJsonParser_UnmarshalTo(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"}`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	var person Person
	err = parser.UnmarshalTo(&person)
	require.NoError(t, err)

	require.Equal(t, "John", person.Name)
	require.Equal(t, 30, person.Age)
	require.Equal(t, "New York", person.City)
}

func TestIncompleteJsonParser_UnmarshalTo_IncompleteJSON(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"city":"New York"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	var person Person
	err = parser.UnmarshalTo(&person)
	require.NoError(t, err)

	require.Equal(t, "John", person.Name)
	require.Equal(t, 30, person.Age)
	require.Equal(t, "New York", person.City)
}

func TestIncompleteJsonParser_UnmarshalTo_NestedStruct(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"John","age":30,"address":{"street":"123 Main St","city":"New York","zipCode":"10001"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	var person PersonWithAddress
	err = parser.UnmarshalTo(&person)
	require.NoError(t, err)

	require.Equal(t, "John", person.Name)
	require.Equal(t, 30, person.Age)
	require.Equal(t, "123 Main St", person.Address.Street)
	require.Equal(t, "New York", person.Address.City)
	require.Equal(t, "10001", person.Address.ZipCode)
}

func TestUnmarshalTo_StaticFunction(t *testing.T) {
	jsonString := `{"name":"Alice","age":25,"city":"Tokyo"}`

	var person Person
	err := UnmarshalTo(jsonString, &person)
	require.NoError(t, err)

	require.Equal(t, "Alice", person.Name)
	require.Equal(t, 25, person.Age)
	require.Equal(t, "Tokyo", person.City)
}

func TestGetObjectsAs_Generics(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"name":"Bob","age":35,"city":"London"}`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	person, err := GetObjectsAs[Person](parser)
	require.NoError(t, err)

	require.Equal(t, "Bob", person.Name)
	require.Equal(t, 35, person.Age)
	require.Equal(t, "London", person.City)
}

func TestParseAs_Generics(t *testing.T) {
	jsonString := `{"name":"Charlie","age":40,"city":"Paris"}`

	person, err := ParseAs[Person](jsonString)
	require.NoError(t, err)

	require.Equal(t, "Charlie", person.Name)
	require.Equal(t, 40, person.Age)
	require.Equal(t, "Paris", person.City)
}

func TestIncompleteJsonParser_UnmarshalTo_BlogPost(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"title": "Go言語の素晴らしさ", "content": "Goは非常にシンプルで効率的なプログラミング言語です。", "tags": ["プログラミング", "Go"], "comments": ["勉強になりました"]}`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	var blog BlogPost
	err = parser.UnmarshalTo(&blog)
	require.NoError(t, err)

	require.Equal(t, "Go言語の素晴らしさ", blog.Title)
	require.Equal(t, "Goは非常にシンプルで効率的なプログラミング言語です。", blog.Content)
	require.Equal(t, []string{"プログラミング", "Go"}, blog.Tags)
	require.Equal(t, []string{"勉強になりました"}, blog.Comments)
}

func TestIncompleteJsonParser_UnmarshalTo_IncompleteBlogPost(t *testing.T) {
	parser := NewIncompleteJsonParser()
	jsonString := `{"title": "日本の四季について", "content": "日本には美しい四季があります。春は桜、夏は祭り、秋は紅葉、冬は雪景色。それぞれの季節に独特の魅力があり、多くの人々を魅了しています。\n\n"`

	err := parser.Write(jsonString)
	require.NoError(t, err)

	var blog BlogPost
	err = parser.UnmarshalTo(&blog)
	require.NoError(t, err)

	require.Equal(t, "日本の四季について", blog.Title)
	require.Equal(t, "日本には美しい四季があります。春は桜、夏は祭り、秋は紅葉、冬は雪景色。それぞれの季節に独特の魅力があり、多くの人々を魅了しています。\n\n", blog.Content)
	require.Empty(t, blog.Tags)
	require.Empty(t, blog.Comments)
}
