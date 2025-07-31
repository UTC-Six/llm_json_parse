package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Parse 通用 JSON 解析器，支持任意对象类型（结构体、map等）
// 兼容标准 JSON、markdown 包裹的 JSON、格式错误的 JSON，并自动修复常见问题。
// T 为目标类型，可以是结构体或 map[string]interface{}。
func Parse[T any](input string) (T, error) {
	var result T

	// 1. 尝试直接解析纯 JSON
	// 适用于标准 JSON 格式，最快捷
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return result, nil
	}

	// 2. 尝试提取 markdown 代码块中的 JSON
	// 兼容 LLM 返回的 markdown 格式，如 ```json ... ```
	extractedJSON := extractJSONFromMarkdown(input)
	if extractedJSON != "" {
		if err := json.Unmarshal([]byte(extractedJSON), &result); err == nil {
			return result, nil
		}
	}

	// 3. 尝试修复格式错误的 JSON
	// 兼容 LLM 返回的非标准 JSON，如缺少引号、单引号、尾逗号等
	_, fixedResult := TryParseJSONObject(input)
	if len(fixedResult) > 0 {
		// 先转为标准 JSON 字符串，再反序列化为目标类型
		jsonBytes, err := json.Marshal(fixedResult)
		if err != nil {
			return result, fmt.Errorf("修复后的 JSON 序列化失败: %v", err)
		}
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return result, fmt.Errorf("修复后的 JSON 解析失败: %v", err)
		}
		return result, nil
	}

	// 4. 所有方式都失败，返回错误
	return result, fmt.Errorf("无法解析任何格式的 JSON")
}

// extractJSONFromMarkdown 从 markdown 文本中提取 JSON 内容
// 支持 ```json ... ``` 和 ``` ... ``` 两种代码块格式
// 只提取第一个代码块，若内容看起来像 JSON 则返回
func extractJSONFromMarkdown(input string) string {
	// 匹配 ```json 和 ``` 之间的内容
	pattern := regexp.MustCompile("```json\\s*\\n([\\s\\S]*?)\\n```")
	matches := pattern.FindStringSubmatch(input)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1])
	}
	// 匹配普通 ``` 代码块
	pattern2 := regexp.MustCompile("```\\s*\\n([\\s\\S]*?)\\n```")
	matches2 := pattern2.FindStringSubmatch(input)
	if len(matches2) >= 2 {
		content := strings.TrimSpace(matches2[1])
		if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
			return content
		}
	}
	return ""
}

// TryParseJSONObject 尝试修复和解析格式错误的 JSON 字符串
// 兼容 LLM 返回的非标准 JSON，自动修复常见问题
// 返回 (修复后的 JSON 字符串, 解析后的 map)
func TryParseJSONObject(input string) (string, map[string]interface{}) {
	var result map[string]interface{}

	// 1. 直接尝试标准解析
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return input, result
	}

	// 2. 尝试提取大括号内容
	pattern := regexp.MustCompile(`\{(.*)\}`)
	matches := pattern.FindStringSubmatch(input)
	if len(matches) >= 2 {
		input = "{" + matches[1] + "}"
	}

	// 3. 清理常见格式问题（如多余括号、转义、换行等）
	input = cleanJSONString(input)

	// 4. 移除 markdown 代码块包裹
	input = removeMarkdownFrame(input)

	// 5. 再次尝试标准解析
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return input, result
	}

	// 6. 最后尝试正则修复（如缺少引号、单引号、尾逗号等）
	jsonInfo := repairJSON(input)
	if err := json.Unmarshal([]byte(jsonInfo), &result); err != nil {
		return jsonInfo, make(map[string]interface{})
	}
	return jsonInfo, result
}

// cleanJSONString 清理常见的 JSON 格式问题
// 包括多余括号、转义字符、换行、回车等
func cleanJSONString(input string) string {
	replacements := map[string]string{
		"{{":   "{",
		"}}":   "}",
		"\"[{": "[{",
		"}]\"": "}]",
		"\\n":  " ",
		"\n":   " ",
		"\r":   "",
		"\\":   " ",
	}
	for old, new := range replacements {
		input = strings.ReplaceAll(input, old, new)
	}
	return strings.TrimSpace(input)
}

// removeMarkdownFrame 移除 markdown 代码块包裹
// 兼容 LLM 返回的 ```json ... ``` 或 ``` ... ``` 格式
func removeMarkdownFrame(input string) string {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "```") {
		input = input[3:]
	}
	if strings.HasPrefix(input, "json") {
		input = input[4:]
	}
	if strings.HasSuffix(input, "```") {
		input = input[:len(input)-3]
	}
	return strings.TrimSpace(input)
}

// repairJSON 尝试修复格式错误的 JSON
// 包括移除尾逗号、补引号、单引号转双引号、字符串值补引号等
func repairJSON(jsonStr string) string {
	// 1. 移除尾逗号
	re := regexp.MustCompile(`,\s*[}\]]`)
	jsonStr = re.ReplaceAllString(jsonStr, "$1")
	// 2. 补全缺失的键名引号
	re = regexp.MustCompile(`(\w+):`)
	jsonStr = re.ReplaceAllString(jsonStr, `"$1":`)
	// 3. 单引号转双引号
	jsonStr = strings.ReplaceAll(jsonStr, "'", "\"")
	// 4. 补全缺失的字符串值引号（不处理数字/布尔/null）
	re = regexp.MustCompile(`:\s*([^"][^,}\]]*[^"\s,}\]])`)
	jsonStr = re.ReplaceAllStringFunc(jsonStr, func(match string) string {
		value := re.FindStringSubmatch(match)[1]
		if value == "true" || value == "false" || value == "null" {
			return match
		}
		if strings.ContainsAny(value, "0123456789") && !strings.ContainsAny(value, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			return match
		}
		return strings.Replace(match, value, `"`+value+`"`, 1)
	})
	return jsonStr
}
