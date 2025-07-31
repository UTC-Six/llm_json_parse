# JSON 解析工具

本仓库包含用于解析和修复 JSON 字符串的工具，最初用 Python 编写，现在也提供 Go 版本。

## Go 版本

Go 实现提供了与 Python 版本相同的功能：

### 函数

#### `Parse(functionString string) (string, map[string]interface{})`

解析函数调用字符串并将其转换为 JSON 格式。

**示例：**
```go
functionString := "tool_call(first_int={'title': 'First Int', 'type': 'integer'}, second_int={'title': 'Second Int', 'type': 'integer'})"
astInfo, jsonResult := Parse(functionString)
```

#### `TryParseJSONObject(input string) (string, map[string]interface{})`

清理和格式化 JSON 字符串，尝试修复格式错误的 JSON。

**示例：**
```go
jsonString := `{"name": "test", "value": 123}`
cleaned, result := TryParseJSONObject(jsonString)
```

### 使用方法

1. 确保安装了 Go 1.21 或更高版本
2. 运行示例：
   ```bash
   go run json_parse.go
   ```

### 功能特性

- **类 AST 解析**：将函数调用字符串转换为 JSON 格式
- **JSON 修复**：尝试修复常见的 JSON 格式问题
- **Markdown 框架移除**：从 JSON 字符串中移除 markdown 代码块
- **类型推断**：自动检测并转换数据类型

### 与 Python 版本的区别

- 使用基于正则表达式的解析而不是 Python 的 `ast` 模块
- 简化的 JSON 修复逻辑（生产环境中可能需要使用更复杂的 JSON 修复库）
- Go 风格的错误处理和返回值
- 除了 Go 标准库外没有外部依赖

## Python 版本

原始的 Python 实现也可用，使用 Python 的 `ast` 模块和 `json_repair` 库提供类似功能。 