package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Parse 解析 JSON 字符串到指定类型
// 支持标准 JSON、markdown 代码块 JSON、以及格式错误的 JSON
func Parse[T any](input string) (*T, error) {
	if input == "" {
		return nil, fmt.Errorf("输入字符串为空")
	}

	// 第一步：提取 JSON 内容
	jsonContent, err := extractJSONContent(input)
	if err != nil {
		return nil, fmt.Errorf("提取 JSON 内容失败: %w", err)
	}

	// 第二步：修复 JSON 格式问题
	fixedJSON, err := fixJSONFormat(jsonContent)
	if err != nil {
		return nil, fmt.Errorf("修复 JSON 格式失败: %w", err)
	}

	// 第三步：解析为对象
	var result T
	if err := json.Unmarshal([]byte(fixedJSON), &result); err != nil {
		return nil, fmt.Errorf("JSON 解析失败: %w", err)
	}

	return &result, nil
}

// extractJSONContent 从输入中提取 JSON 内容
// 支持标准 JSON 和 markdown 代码块格式
func extractJSONContent(input string) (string, error) {
	input = strings.TrimSpace(input)

	// 尝试匹配 markdown 代码块
	markdownPattern := regexp.MustCompile("```(?:json)?\\s*\\n?(?s:(.*?))\\n?```")
	matches := markdownPattern.FindStringSubmatch(input)

	if len(matches) > 1 {
		// 找到 markdown 代码块，返回其中的内容
		return strings.TrimSpace(matches[1]), nil
	}

	// 如果没有找到 markdown 代码块，检查是否是纯 JSON
	if strings.HasPrefix(input, "{") || strings.HasPrefix(input, "[") {
		return input, nil
	}

	return "", fmt.Errorf("未找到有效的 JSON 内容")
}

// fixJSONFormat 修复 JSON 格式问题
func fixJSONFormat(jsonStr string) (string, error) {
	if jsonStr == "" {
		return "", fmt.Errorf("JSON 字符串为空")
	}

	// 先尝试直接解析，如果成功就直接返回
	var temp interface{}
	if err := json.Unmarshal([]byte(jsonStr), &temp); err == nil {
		return jsonStr, nil
	}

	// 如果解析失败，进行修复
	fixed := jsonStr

	// 1. 转换中文引号
	fixed = convertChineseQuotes(fixed)

	// 2. 修复缺失的引号
	fixed = fixMissingQuotes(fixed)

	// 3. 修复缺失的括号
	fixed = fixMissingBrackets(fixed)

	// 4. 修复尾随逗号
	fixed = fixTrailingCommas(fixed)

	// 5. 修复未转义的引号
	fixed = fixUnescapedQuotes(fixed)

	// 6. 修复单引号
	fixed = fixSingleQuotes(fixed)

	// 再次尝试解析
	if err := json.Unmarshal([]byte(fixed), &temp); err != nil {
		return "", fmt.Errorf("修复后仍无法解析 JSON: %w", err)
	}

	return fixed, nil
}

// convertChineseQuotes 将中文双引号转换为英文双引号
func convertChineseQuotes(input string) string {
	// 中文双引号：""（左双引号）和 ""（右双引号）
	input = strings.ReplaceAll(input, "\u201c", `"`)
	input = strings.ReplaceAll(input, "\u201d", `"`)
	return input
}

// fixMissingQuotes 修复缺失的引号
func fixMissingQuotes(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 修复 key 缺失引号的情况：key: value -> "key": value
		colonIndex := strings.Index(line, ":")
		if colonIndex > 0 {
			key := strings.TrimSpace(line[:colonIndex])
			value := strings.TrimSpace(line[colonIndex+1:])

			// 如果 key 没有引号且不是数字或布尔值
			if !strings.HasPrefix(key, `"`) && !strings.HasSuffix(key, `"`) {
				// 检查是否是有效的 key（不包含特殊字符）
				if isValidKey(key) {
					key = `"` + key + `"`
				}
			}

			// 如果 value 是字符串但没有引号
			if value != "" && !strings.HasPrefix(value, `"`) && !strings.HasSuffix(value, `"`) {
				// 检查是否是字符串值（不是数字、布尔值、null、对象、数组）
				if isStringValue(value) {
					value = `"` + value + `"`
				}
			}

			line = key + ": " + value
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// fixMissingBrackets 修复缺失的括号
func fixMissingBrackets(input string) string {
	// 统计括号数量
	openBraces := strings.Count(input, "{")
	closeBraces := strings.Count(input, "}")
	openBrackets := strings.Count(input, "[")
	closeBrackets := strings.Count(input, "]")

	// 修复缺失的大括号
	if openBraces > closeBraces {
		input += strings.Repeat("}", openBraces-closeBraces)
	}

	// 修复缺失的方括号
	if openBrackets > closeBrackets {
		input += strings.Repeat("]", openBrackets-closeBrackets)
	}

	return input
}

// fixTrailingCommas 修复尾随逗号
func fixTrailingCommas(input string) string {
	// 移除对象和数组末尾的逗号
	input = regexp.MustCompile(`,\s*([}\]])`).ReplaceAllString(input, "$1")
	return input
}

// fixUnescapedQuotes 修复未转义的引号
func fixUnescapedQuotes(input string) string {
	var result strings.Builder
	inString := false
	escaped := false
	pos := 0

	for pos < len(input) {
		char := input[pos]

		if escaped {
			result.WriteByte(char)
			escaped = false
			pos++
			continue
		}

		if char == '\\' {
			result.WriteByte(char)
			escaped = true
			pos++
			continue
		}

		if char == '"' {
			if !inString {
				inString = true
				result.WriteByte(char)
			} else {
				if isStringEndPosition(input, pos) {
					inString = false
					result.WriteByte(char)
				} else {
					result.WriteString(`\"`)
				}
			}
			pos++
			continue
		}

		result.WriteByte(char)
		pos++
	}

	return result.String()
}

// fixSingleQuotes 修复单引号
func fixSingleQuotes(input string) string {
	// 将单引号替换为双引号，但需要处理转义
	var result strings.Builder
	inString := false
	escaped := false
	pos := 0

	for pos < len(input) {
		char := input[pos]

		if escaped {
			result.WriteByte(char)
			escaped = false
			pos++
			continue
		}

		if char == '\\' {
			result.WriteByte(char)
			escaped = true
			pos++
			continue
		}

		if char == '\'' {
			if !inString {
				inString = true
				result.WriteByte('"')
			} else {
				if isStringEndPosition(input, pos) {
					inString = false
					result.WriteByte('"')
				} else {
					result.WriteString(`\"`)
				}
			}
			pos++
			continue
		}

		result.WriteByte(char)
		pos++
	}

	return result.String()
}

// isValidKey 检查是否是有效的 key
func isValidKey(key string) bool {
	// 简单的 key 验证：不包含特殊字符且不是数字开头
	if len(key) == 0 {
		return false
	}

	// 检查是否以数字开头
	if key[0] >= '0' && key[0] <= '9' {
		return false
	}

	// 检查是否包含特殊字符（除了字母、数字、下划线）
	for _, char := range key {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_') {
			return false
		}
	}

	return true
}

// isStringValue 检查是否是字符串值
func isStringValue(value string) bool {
	value = strings.TrimSpace(value)

	// 如果包含中文字符，一定是字符串
	for _, char := range value {
		if char >= 0x4e00 && char <= 0x9fff {
			return true
		}
	}

	// 检查是否是数字
	if regexp.MustCompile(`^-?\d+(\.\d+)?$`).MatchString(value) {
		return false
	}
	// 检查是否是布尔值
	if value == "true" || value == "false" {
		return false
	}
	// 检查是否是 null
	if value == "null" {
		return false
	}
	// 检查是否是对象或数组
	if strings.HasPrefix(value, "{") || strings.HasPrefix(value, "[") {
		return false
	}
	return true
}

// isStringEndPosition 判断当前位置的引号是否是字符串结束
func isStringEndPosition(input string, pos int) bool {
	pos++
	// 跳过空白字符
	for pos < len(input) {
		char := input[pos]
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			pos++
			continue
		}
		break
	}
	if pos >= len(input) {
		return true
	}
	nextChar := input[pos]
	return nextChar == ',' || nextChar == ':' || nextChar == '}' || nextChar == ']'
}
