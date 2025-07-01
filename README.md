# Incomplete JSON Parser - Go Port

This is a Go port of the [1000ship/incomplete-json-parser](https://github.com/1000ship/incomplete-json-parser)

## Installation

```bash
go get github.com/kiokuless/incomplete-json-parser-go
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    incompletejson "github.com/kiokuless/incomplete-json-parser-go"
)

func main() {
    // Create a new parser
    parser := incompletejson.NewIncompleteJsonParser()
    
    // Parse incomplete JSON
    err := parser.Write(`{"name":"John","age":30,"city":"New`)
    if err != nil {
        log.Fatal(err)
    }
    
    result, err := parser.GetObjects()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %+v\n", result)
    // Output: Result: map[age:30 city:New name:John]
}
```

### Type-Safe Parsing with UnmarshalTo

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
    City string `json:"city"`
}

func main() {
    parser := incompletejson.NewIncompleteJsonParser()
    parser.Write(`{"name":"John","age":30,"city":"New York"}`)
    
    var person Person
    err := parser.UnmarshalTo(&person)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Name: %s, Age: %d, City: %s\n", person.Name, person.Age, person.City)
    // Output: Name: John, Age: 30, City: New York
}
```

### Static Functions

```go
// Basic parsing
result, err := incompletejson.Parse(`{"name":"John","age":30}`)

// Type-safe parsing
var person Person
err := incompletejson.UnmarshalTo(`{"name":"Alice","age":25}`, &person)
```

### Generics Support (Go 1.18+)

```go
// Using generics for type-safe parsing
person, err := incompletejson.ParseAs[Person](`{"name":"Bob","age":35}`)
if err != nil {
    log.Fatal(err)
}

// Or with parser instance
parser := incompletejson.NewIncompleteJsonParser()
parser.Write(`{"name":"Charlie","age":40}`)
person, err := incompletejson.GetObjectsAs[Person](parser)
```

### Advanced Options

```go
// Ignore extra characters after valid JSON
parser := incompletejson.NewIncompleteJsonParser(
    incompletejson.WithIgnoreExtraCharacters(true),
)

// Parse JSON with trailing text
parser.Write(`{"message":"Hello"}\n\nExtra text here`)
result, _ := parser.GetObjects() // Works without error
```

## Testing

Run the tests:

```bash
go test -v
```

## Build

```bash
go build
```

## Features

### Core Features
- Parse incomplete JSON objects and arrays
- Parse incomplete JSON strings
- Handle null values with different lengths
- Support for nested objects and arrays
- Streaming parser that can handle multiple chunks

### Type Safety Features
- **UnmarshalTo**: Type-safe parsing with struct mapping
- **Generics Support**: Modern Go generics for compile-time type safety
- **JSON Tags**: Full support for standard `json:` tags
- **Static Functions**: Convenient one-line parsing

### Advanced Options
- **WithIgnoreExtraCharacters**: Option to ignore text after valid JSON
- **Functional Options**: Clean API for parser configuration

## API Reference

### Constructor
```go
// Basic parser
parser := NewIncompleteJsonParser()

// Parser with options
parser := NewIncompleteJsonParser(WithIgnoreExtraCharacters(true))
```

### Instance Methods
```go
// Write JSON data (can be called multiple times)
err := parser.Write(jsonString)

// Get parsed result as interface{}
result, err := parser.GetObjects()

// Type-safe parsing
var target MyStruct
err := parser.UnmarshalTo(&target)

// Reset parser state
parser.Reset()
```

### Static Functions
```go
// Basic parsing
result, err := Parse(jsonString)

// Type-safe parsing
var target MyStruct
err := UnmarshalTo(jsonString, &target)

// Generics (Go 1.18+)
result, err := ParseAs[MyStruct](jsonString)
target, err := GetObjectsAs[MyStruct](parser)
```

## Error Handling

The parser handles various incomplete JSON scenarios gracefully:

```go
// Missing closing braces
`{"name":"John","age":30` → map[name:John age:30]

// Incomplete strings
`{"message":"Hello, world` → map[message:"Hello, world"]

// Incomplete arrays
`["apple","banana","orange"` → ["apple","banana","orange"]

// Trailing commas
`{"a":1,"b":2,` → map[a:1 b:2]
```
