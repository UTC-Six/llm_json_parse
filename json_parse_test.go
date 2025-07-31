package main

import (
	"reflect"
	"testing"
)

// TestParse 测试 Parse 函数的基本功能
func TestParse(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult map[string]interface{}
	}{
		{
			name:  "基本函数调用",
			input: "tool_call(first_int={'title': 'First Int', 'type': 'integer'})",
			expectedResult: map[string]interface{}{
				"first_int": map[string]interface{}{
					"title": "First Int",
					"type":  "integer",
				},
			},
		},
		{
			name:  "多个参数",
			input: "function(a=1, b=2, c='hello')",
			expectedResult: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "hello",
			},
		},
		{
			name:  "嵌套对象参数",
			input: "test(data={'name': 'test', 'value': 123, 'active': true})",
			expectedResult: map[string]interface{}{
				"data": map[string]interface{}{
					"name":   "test",
					"value":  123,
					"active": true,
				},
			},
		},
		{
			name:           "空输入",
			input:          "",
			expectedResult: map[string]interface{}{},
		},
		{
			name:           "无效格式",
			input:          "invalid_format",
			expectedResult: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, result := Parse(tt.input)
			
			// 检查结果
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Parse() result = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

// TestParseArguments 测试 parseArguments 函数
func TestParseArguments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected map[string]interface{}
	}{
		{
			name:  "简单参数",
			input: "a=1, b=2",
			expected: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
		},
		{
			name:  "字符串参数",
			input: "name='test', value='hello'",
			expected: map[string]interface{}{
				"name":  "test",
				"value": "hello",
			},
		},
		{
			name:  "布尔值参数",
			input: "active=true, enabled=false",
			expected: map[string]interface{}{
				"active":  true,
				"enabled": false,
			},
		},
		{
			name:  "嵌套对象",
			input: "config={'key': 'value', 'number': 42}",
			expected: map[string]interface{}{
				"config": map[string]interface{}{
					"key":    "value",
					"number": 42,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArguments(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseArguments() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSplitArgs 测试 splitArgs 函数
func TestSplitArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "简单分割",
			input:    "a=1, b=2, c=3",
			expected: []string{"a=1", "b=2", "c=3"},
		},
		{
			name:     "嵌套对象",
			input:    "a={x:1,y:2}, b=[1,2,3]",
			expected: []string{"a={x:1,y:2}", "b=[1 2 3]"},
		},
		{
			name:     "复杂嵌套",
			input:    "config={data:{name:'test'}}, enabled=true",
			expected: []string{"config={data:{name:'test'}}", "enabled=true"},
		},
		{
			name:     "空输入",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := splitArgs(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("splitArgs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestParseValue 测试 parseValue 函数
func TestParseValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "整数",
			input:    "123",
			expected: 123,
		},
		{
			name:     "浮点数",
			input:    "3.14",
			expected: 3.14,
		},
		{
			name:     "布尔值 true",
			input:    "true",
			expected: true,
		},
		{
			name:     "布尔值 false",
			input:    "false",
			expected: false,
		},
		{
			name:     "字符串",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "带引号的字符串",
			input:    "'test'",
			expected: "test",
		},
		{
			name:     "JSON 对象",
			input:    "{'name': 'test', 'value': 123}",
			expected: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseValue(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("parseValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestTryParseJSONObject 测试 TryParseJSONObject 函数
func TestTryParseJSONObject(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedResult map[string]interface{}
		shouldSucceed  bool
	}{
		{
			name:  "有效 JSON",
			input: `{"name": "test", "value": 123}`,
			expectedResult: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			shouldSucceed: true,
		},
		{
			name:  "带 Markdown 的 JSON",
			input: "```json\n{\"name\": \"test\"}\n```",
			expectedResult: map[string]interface{}{
				"name": "test",
			},
			shouldSucceed: true,
		},
		{
			name:  "格式错误的 JSON",
			input: "{name: 'test', value: 123}",
			expectedResult: map[string]interface{}{
				"name":  "test",
				"value": 123,
			},
			shouldSucceed: true,
		},
		{
			name:           "无效 JSON",
			input:          "invalid json",
			expectedResult: map[string]interface{}{},
			shouldSucceed:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, result := TryParseJSONObject(tt.input)
			
			if tt.shouldSucceed {
				if !reflect.DeepEqual(result, tt.expectedResult) {
					t.Errorf("TryParseJSONObject() result = %v, want %v", result, tt.expectedResult)
				}
			} else {
				if len(result) != 0 {
					t.Errorf("TryParseJSONObject() should return empty result for invalid input, got %v", result)
				}
			}
		})
	}
}

// TestCleanJSONString 测试 cleanJSONString 函数
func TestCleanJSONString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "双花括号",
			input:    "{{name: 'test'}}",
			expected: "{name: 'test'}",
		},
		{
			name:     "转义字符",
			input:    "{\"name\": \"test\\nvalue\"}",
			expected: "{\"name\": \"test value\"}",
		},
		{
			name:     "数组标记",
			input:    "\"[{data: 1}]\"",
			expected: "[{data: 1}]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cleanJSONString(tt.input)
			if result != tt.expected {
				t.Errorf("cleanJSONString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestRemoveMarkdownFrame 测试 removeMarkdownFrame 函数
func TestRemoveMarkdownFrame(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "基本 Markdown",
			input:    "```json\n{\"name\": \"test\"}\n```",
			expected: "{\"name\": \"test\"}",
		},
		{
			name:     "只有开始标记",
			input:    "```json\n{\"name\": \"test\"}",
			expected: "{\"name\": \"test\"}",
		},
		{
			name:     "没有 Markdown",
			input:    "{\"name\": \"test\"}",
			expected: "{\"name\": \"test\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := removeMarkdownFrame(tt.input)
			if result != tt.expected {
				t.Errorf("removeMarkdownFrame() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestRepairJSON 测试 repairJSON 函数
func TestRepairJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "尾随逗号",
			input:    "{name: 'test', value: 123,}",
			expected: "{\"name\": \"test\", \"value\": 123}",
		},
		{
			name:     "缺失引号",
			input:    "{name: 'test', value: 123}",
			expected: "{\"name\": \"test\", \"value\": 123}",
		},
		{
			name:     "单引号转双引号",
			input:    "{'name': 'test'}",
			expected: "{\"name\": \"test\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := repairJSON(tt.input)
			if result != tt.expected {
				t.Errorf("repairJSON() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// BenchmarkParse 性能测试
func BenchmarkParse(b *testing.B) {
	input := "tool_call(first_int={'title': 'First Int', 'type': 'integer'}, second_int={'title': 'Second Int', 'type': 'integer'})"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Parse(input)
	}
}

// BenchmarkTryParseJSONObject 性能测试
func BenchmarkTryParseJSONObject(b *testing.B) {
	input := `{"name": "test", "value": 123, "active": true}`
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		TryParseJSONObject(input)
	}
}

// TestIntegration 集成测试
func TestIntegration(t *testing.T) {
	// 测试完整的解析流程
	functionString := "api_call(user={'id': 123, 'name': 'test'}, settings={'enabled': true, 'timeout': 30})"
	
	_, jsonResult := Parse(functionString)
	
	// 验证结果包含预期的键
	expectedKeys := []string{"user", "settings"}
	for _, key := range expectedKeys {
		if _, exists := jsonResult[key]; !exists {
			t.Errorf("Integration test: Result should contain key '%s'", key)
		}
	}
	
	// 验证嵌套对象
	if user, ok := jsonResult["user"].(map[string]interface{}); ok {
		if user["id"] != float64(123) { // JSON 解析后数字是 float64
			t.Error("Integration test: User ID should be 123")
		}
	} else {
		t.Error("Integration test: User should be a map")
	}
}
