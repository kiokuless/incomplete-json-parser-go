# Incomplete JSON Parser - Go Port

This is a Go port of the [1000ship/incomplete-json-parser](https://github.com/1000ship/incomplete-json-parser)

## Installation

```bash
go get github.com/kiokuless/incomplete-json-parser-go
```

## Usage

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

## Static Parse Method

```go
result, err := incompletejson.Parse(`{"name":"John","age":30}`)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Result: %+v\n", result)
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

- Parse incomplete JSON objects
- Parse incomplete JSON arrays  
- Parse incomplete JSON strings
- Handle null values with different lengths
- Support for nested objects and arrays
- Streaming parser that can handle multiple chunks
