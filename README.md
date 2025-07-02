# Incomplete JSON Parser - Go Port

[![CI](https://github.com/kiokuless/incomplete-json-parser-go/workflows/CI/badge.svg)](https://github.com/kiokuless/incomplete-json-parser-go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/kiokuless/incomplete-json-parser-go)](https://goreportcard.com/report/github.com/kiokuless/incomplete-json-parser-go)
[![codecov](https://codecov.io/gh/kiokuless/incomplete-json-parser-go/branch/main/graph/badge.svg)](https://codecov.io/gh/kiokuless/incomplete-json-parser-go)
[![Go Reference](https://pkg.go.dev/badge/github.com/kiokuless/incomplete-json-parser-go.svg)](https://pkg.go.dev/github.com/kiokuless/incomplete-json-parser-go)
[![Latest Release](https://img.shields.io/github/v/release/kiokuless/incomplete-json-parser-go)](https://github.com/kiokuless/incomplete-json-parser-go/releases)
[![License](https://img.shields.io/github/license/kiokuless/incomplete-json-parser-go)](LICENSE)

This is a Go port of the [1000ship/incomplete-json-parser](https://github.com/1000ship/incomplete-json-parser) TypeScript library.

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

// Allow unescaped newlines in JSON strings
parser := incompletejson.NewIncompleteJsonParser(
    incompletejson.WithAllowUnescapedNewlines(true),
)

// Parse JSON with literal newlines in strings
parser.Write(`{"text": "Hello
World"}`)
result, _ := parser.GetObjects() // result: map[text:Hello\nWorld]

// Using options with static functions
result, err := incompletejson.Parse(`{"text": "Hello
World"}`, incompletejson.WithAllowUnescapedNewlines(true))

var data MyStruct
err := incompletejson.UnmarshalTo(`{"text": "Hello
World"}`, &data, incompletejson.WithAllowUnescapedNewlines(true))

result, err := incompletejson.ParseAs[MyStruct](`{"text": "Hello
World"}`, incompletejson.WithAllowUnescapedNewlines(true))

// Validate required fields (non-omitempty)
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Email string `json:"email,omitempty"` // optional
}

var user User
err := incompletejson.UnmarshalTo(`{"id": 1}`, &user, incompletejson.WithRequiredFields(true))
// Error: missing required fields: name

// With ParseAs
user, err := incompletejson.ParseAs[User](`{"id": 1, "name": "John"}`, incompletejson.WithRequiredFields(true))
// Success: email is optional (omitempty)
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
- **WithAllowUnescapedNewlines**: Option to allow unescaped newlines in JSON strings
- **WithRequiredFields**: Option to validate that all non-omitempty fields are present
- **Functional Options**: Clean API for parser configuration

## API Reference

### Constructor
```go
// Basic parser
parser := NewIncompleteJsonParser()

// Parser with options
parser := NewIncompleteJsonParser(
    WithIgnoreExtraCharacters(true),
    WithAllowUnescapedNewlines(true),
    WithRequiredFields(true),
)
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

// Basic parsing with options
result, err := Parse(jsonString, WithAllowUnescapedNewlines(true))

// Type-safe parsing
var target MyStruct
err := UnmarshalTo(jsonString, &target)

// Type-safe parsing with options
err := UnmarshalTo(jsonString, &target, WithAllowUnescapedNewlines(true), WithRequiredFields(true))

// Generics (Go 1.18+)
result, err := ParseAs[MyStruct](jsonString)
result, err := ParseAs[MyStruct](jsonString, WithRequiredFields(true))
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

## Development

This project uses [mise](https://mise.jdx.dev/) for development tools and task management.

### Setup

```bash
# Install mise (if not already installed)
curl https://mise.run | sh

# Install development tools
mise install

# Install pre-commit hooks
mise run precommit-install
```

### Common Tasks

```bash
# Run tests
mise run test

# Run tests with coverage
mise run test-coverage

# Run linter
mise run lint

# Format code
mise run fmt

# Run all checks
mise run check

# Check next version
mise run version-check

# Release (patch/minor/major)
mise run release-patch
mise run release-minor
mise run release-major
```

### Release Process

1. Make sure you're on the main branch with a clean working directory
2. Check what the next version would be: `mise run version-check`
3. Release: `mise run release-patch` (or `release-minor`/`release-major`)
4. The script will run tests, create a tag, and push to GitHub
5. GitHub Actions will automatically create the release

### Available Commands

Run `mise tasks` to see all available commands.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
