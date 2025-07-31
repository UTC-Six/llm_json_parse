# JSON Parse Utilities

This repository contains utilities for parsing and repairing JSON strings, originally written in Python and now also available in Go.

## Go Version

The Go implementation provides the same functionality as the Python version:

### Functions

#### `Parse(functionString string) (string, map[string]interface{})`

Parses a function call string and converts it to JSON format.

**Example:**
```go
functionString := "tool_call(first_int={'title': 'First Int', 'type': 'integer'}, second_int={'title': 'Second Int', 'type': 'integer'})"
astInfo, jsonResult := Parse(functionString)
```

#### `TryParseJSONObject(input string) (string, map[string]interface{})`

Cleans and formats JSON strings, attempting to repair malformed JSON.

**Example:**
```go
jsonString := `{"name": "test", "value": 123}`
cleaned, result := TryParseJSONObject(jsonString)
```

### Usage

1. Make sure you have Go 1.21 or later installed
2. Run the example:
   ```bash
   go run json_parse.go
   ```

### Features

- **AST-like parsing**: Converts function call strings to JSON format
- **JSON repair**: Attempts to fix common JSON formatting issues
- **Markdown frame removal**: Removes markdown code blocks from JSON strings
- **Type inference**: Automatically detects and converts data types

### Differences from Python Version

- Uses regex-based parsing instead of Python's `ast` module
- Simplified JSON repair logic (you may want to use a more sophisticated JSON repair library for production use)
- Go-style error handling and return values
- No external dependencies beyond Go standard library

## Python Version

The original Python implementation is also available and provides similar functionality using Python's `ast` module and the `json_repair` library.