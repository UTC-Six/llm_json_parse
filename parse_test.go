package main

import (
	"os"
	"testing"
)

// TestStruct 用于测试通用 JSON 解析的结构体
// 字段类型覆盖 string、int、bool
// 用于验证泛型解析的类型安全
type TestStruct struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	IsValid bool   `json:"is_valid"`
}

// ThinkResponse 用于测试 LLM markdown JSON 解析的结构体
type ThinkResponse struct {
	Thinking string `json:"thinking"`
	Chat     string `json:"chat"`
}

// TestParsePureJSON 测试标准 JSON 字符串的解析能力
// 验证 Parse 能否直接解析标准 JSON 到结构体
func TestParsePureJSON(t *testing.T) {
	pureJSON := `{
		"name": "张三",
		"age": 25,
		"is_valid": true
	}`

	result, err := Parse[TestStruct](pureJSON)
	if err != nil {
		t.Errorf("解析失败: %v", err)
		return
	}

	if result.Name != "张三" {
		t.Errorf("期望 Name 为 '张三'，实际为 '%s'", result.Name)
	}
	if result.Age != 25 {
		t.Errorf("期望 Age 为 25，实际为 %d", result.Age)
	}
	if !result.IsValid {
		t.Errorf("期望 IsValid 为 true，实际为 %t", result.IsValid)
	}
}

// TestParseMap 测试标准 JSON 解析为 map
// 验证 Parse 能否直接解析为 map[string]interface{}
func TestParseMap(t *testing.T) {
	pureJSON := `{
		"name": "李四",
		"age": 30,
		"is_valid": false
	}`

	result, err := Parse[map[string]interface{}](pureJSON)
	if err != nil {
		t.Errorf("解析失败: %v", err)
		return
	}

	if result["name"] != "李四" {
		t.Errorf("期望 name 为 '李四'，实际为 '%v'", result["name"])
	}
	if result["age"].(float64) != 30 {
		t.Errorf("期望 age 为 30，实际为 %v", result["age"])
	}
	if result["is_valid"] != false {
		t.Errorf("期望 is_valid 为 false，实际为 %v", result["is_valid"])
	}
}

// TestParseMarkdownJSON 测试 markdown 代码块 JSON 的解析能力
// 验证 Parse 能否自动提取并解析 markdown 格式的 JSON
func TestParseMarkdownJSON(t *testing.T) {
	markdownJSON := "这是一个响应：\n\n```json\n{\n    \"thinking\": \"分析文本内容\",\n    \"chat\": \"<span>处理结果</span>\"\n}\n```\n\n结束。"

	result, err := Parse[ThinkResponse](markdownJSON)
	if err != nil {
		t.Errorf("解析失败: %v", err)
		return
	}

	if result.Thinking != "分析文本内容" {
		t.Errorf("期望 Thinking 为 '分析文本内容'，实际为 '%s'", result.Thinking)
	}
	if result.Chat != "<span>处理结果</span>" {
		t.Errorf("期望 Chat 为 '<span>处理结果</span>'，实际为 '%s'", result.Chat)
	}
}

// TestParseMalformedJSON 测试格式错误 JSON 的修复与解析能力
// 验证 Parse 能否自动修复缺少引号、单引号、尾逗号等问题
func TestParseMalformedJSON(t *testing.T) {
	malformedJSON := `{
		name: 张三,
		age: 25,
		is_valid: true
	}`

	result, err := Parse[TestStruct](malformedJSON)
	if err != nil {
		t.Errorf("解析失败: %v", err)
		return
	}

	if result.Name != "张三" {
		t.Errorf("期望 Name 为 '张三'，实际为 '%s'", result.Name)
	}
	if result.Age != 25 {
		t.Errorf("期望 Age 为 25，实际为 %d", result.Age)
	}
	if !result.IsValid {
		t.Errorf("期望 IsValid 为 true，实际为 %t", result.IsValid)
	}
}

// TestParseFromFile 测试从文件读取并解析
// 验证 Parse 能否处理文件中的标准或 markdown JSON
func TestParseFromFile(t *testing.T) {
	content, err := os.ReadFile("json.txt")
	if err != nil {
		t.Skipf("跳过文件测试，无法读取 json.txt: %v", err)
		return
	}

	result, err := Parse[ThinkResponse](string(content))
	if err != nil {
		t.Errorf("解析失败: %v", err)
		return
	}

	// 验证结果不为空
	if result.Thinking == "" {
		t.Errorf("Thinking 字段为空")
	}
	if result.Chat == "" {
		t.Errorf("Chat 字段为空")
	}

	// 验证包含预期的内容
	if !contains(result.Thinking, "军帽下的视线锐利如刀锋") {
		t.Errorf("Thinking 中应包含 '军帽下的视线锐利如刀锋'，实际为: %s", result.Thinking)
	}
	if !contains(result.Chat, "background-color: lightgreen") {
		t.Errorf("Chat 中应包含 'background-color: lightgreen'，实际为: %s", result.Chat)
	}
}

// TestParseInvalidJSON 测试无效 JSON 的报错能力
// 验证 Parse 对于完全无效的 JSON 能否正确报错
func TestParseInvalidJSON(t *testing.T) {
	invalidJSON := `{invalid json content`

	_, err := Parse[TestStruct](invalidJSON)
	if err == nil {
		t.Errorf("期望解析失败，但解析成功了")
	}
}

// TestParseEmptyString 测试空字符串的报错能力
// 验证 Parse 对于空输入能否正确报错
func TestParseEmptyString(t *testing.T) {
	_, err := Parse[TestStruct]("")
	if err == nil {
		t.Errorf("期望解析失败，但解析成功了")
	}
}

// contains 辅助函数，检查字符串是否包含子字符串
// 用于验证解析结果中是否包含预期内容
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			func() bool {
				for i := 0; i <= len(s)-len(substr); i++ {
					if s[i:i+len(substr)] == substr {
						return true
					}
				}
				return false
			}())))
}

// BenchmarkParse 性能测试，验证大批量解析时的性能
func BenchmarkParse(b *testing.B) {
	pureJSON := `{
		"name": "张三",
		"age": 25,
		"is_valid": true
	}`

	for i := 0; i < b.N; i++ {
		_, err := Parse[TestStruct](pureJSON)
		if err != nil {
			b.Errorf("解析失败: %v", err)
		}
	}
}

// BenchmarkParseMarkdown 性能测试 - Markdown JSON
// 验证 markdown 格式 JSON 的批量解析性能
func BenchmarkParseMarkdown(b *testing.B) {
	markdownJSON := "```json\n{\n    \"thinking\": \"分析文本内容\",\n    \"chat\": \"<span>处理结果</span>\"\n}\n```"

	for i := 0; i < b.N; i++ {
		_, err := Parse[ThinkResponse](markdownJSON)
		if err != nil {
			b.Errorf("解析失败: %v", err)
		}
	}
}
