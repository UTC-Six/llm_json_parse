# Universal JSON Parser

A universal JSON parser for Go that can handle various JSON formats including:
- Standard JSON
- Markdown-wrapped JSON
- Malformed JSON with automatic repair

## Features

- **Generic Support**: Use with any struct or `map[string]interface{}`
- **Markdown Extraction**: Automatically extracts JSON from markdown code blocks
- **JSON Repair**: Fixes common JSON formatting issues
- **Multiple Formats**: Handles pure JSON, markdown JSON, and malformed JSON
- **Type Safety**: Full type safety with Go generics

## Quick Start

```go
package main

import (
    "fmt"
    "your-module/parser"
)

// Define your struct
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // Parse standard JSON
    jsonStr := `{"name": "Alice", "age": 30}`
    user, err := Parse[User](jsonStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("User: %+v\n", user)
    
    // Parse markdown JSON
    markdownStr := "```json\n{\"name\": \"Bob\", \"age\": 25}\n```"
    user2, err := Parse[User](markdownStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("User: %+v\n", user2)
    
    // Parse malformed JSON
    malformedStr := `{name: "Charlie", age: 35}`
    user3, err := Parse[User](malformedStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("User: %+v\n", user3)
}
```

## API Reference

### Parse[T]

Parses JSON into any type T.

```go
func Parse[T any](input string) (T, error)
```

**Parameters:**
- `input`: JSON string (can be pure JSON, markdown JSON, or malformed JSON)

**Returns:**
- Parsed object of type T
- Error if parsing fails

## Supported Formats

1. **Standard JSON**
   ```json
   {"name": "value", "number": 123}
   ```

2. **Markdown JSON**
   ```markdown
   ```json
   {"name": "value", "number": 123}
   ```
   ```

3. **Malformed JSON** (automatically repaired)
   ```json
   {name: value, number: 123}
   ```

## Examples

See `parser.go` for complete examples including:
- Struct parsing
- Map parsing
- Markdown extraction
- JSON repair
- File reading

## Installation

```bash
go get github.com/your-username/universal-json-parser
```

## License

MIT License