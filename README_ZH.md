# 通用 JSON 解析器

一个功能强大的 Go JSON 解析器，能够处理各种 JSON 格式：
- 标准 JSON
- Markdown 包装的 JSON
- 格式错误的 JSON（自动修复）

## 特性

- **泛型支持**: 可与任何结构体或 `map[string]interface{}` 一起使用
- **Markdown 提取**: 自动从 markdown 代码块中提取 JSON
- **JSON 修复**: 修复常见的 JSON 格式问题
- **多格式支持**: 处理纯 JSON、markdown JSON 和格式错误的 JSON
- **类型安全**: 使用 Go 泛型提供完整的类型安全

## 快速开始

```go
package main

import (
    "fmt"
    "your-module/parser"
)

// 定义你的结构体
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // 解析标准 JSON
    jsonStr := `{"name": "张三", "age": 30}`
    user, err := Parse[User](jsonStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户: %+v\n", user)
    
    // 解析 markdown JSON
    markdownStr := "```json\n{\"name\": \"李四\", \"age\": 25}\n```"
    user2, err := Parse[User](markdownStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户: %+v\n", user2)
    
    // 解析格式错误的 JSON
    malformedStr := `{name: "王五", age: 35}`
    user3, err := Parse[User](malformedStr)
    if err != nil {
        panic(err)
    }
    fmt.Printf("用户: %+v\n", user3)
}
```

## API 参考

### Parse[T]

将 JSON 解析为任意类型 T。

```go
func Parse[T any](input string) (T, error)
```

**参数:**
- `input`: JSON 字符串（可以是纯 JSON、markdown JSON 或格式错误的 JSON）

**返回:**
- 类型为 T 的解析对象
- 如果解析失败则返回错误

## 支持的格式

1. **标准 JSON**
   ```json
   {"name": "value", "number": 123}
   ```

2. **Markdown JSON**
   ```markdown
   ```json
   {"name": "value", "number": 123}
   ```
   ```

3. **格式错误的 JSON**（自动修复）
   ```json
   {name: value, number: 123}
   ```

## 示例

查看 `parser.go` 获取完整示例，包括：
- 结构体解析
- Map 解析
- Markdown 提取
- JSON 修复
- 文件读取

## 安装

```bash
go get github.com/your-username/universal-json-parser
```

## 许可证

MIT 许可证 