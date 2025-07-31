package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// Parse 解析函数字符串并将其转换为 JSON 格式
// 这个函数模拟了 Python 的 ast 模块功能，用于解析函数调用字符串
// 示例输入: "tool_call(first_int={'title': 'First Int', 'type': 'integer'}, second_int={'title': 'Second Int', 'type': 'integer'})"
// 返回值: (解析信息字符串, JSON 结果映射)
func Parse(functionString string) (string, map[string]interface{}) {
	astInfo := ""                              // 存储解析过程的详细信息
	jsonResult := make(map[string]interface{}) // 存储最终的 JSON 结果

	// 清理输入字符串，去除首尾空白字符
	functionString = strings.TrimSpace(functionString)

	// 使用正则表达式解析函数调用
	// 正则模式: (\w+)\(([^)]*)\)
	// - (\w+): 匹配函数名（一个或多个字母数字字符）
	// - \(: 匹配左括号
	// - ([^)]*): 匹配括号内的所有内容（除了右括号）
	// - \): 匹配右括号
	// 这是对 Python ast 模块的简化版本
	re := regexp.MustCompile(`(\w+)\(([^)]*)\)`)
	matches := re.FindStringSubmatch(functionString)

	// 检查是否成功匹配（至少需要3个元素：完整匹配、函数名、参数）
	if len(matches) >= 3 {
		functionName := matches[1] // 提取函数名
		argsString := matches[2]   // 提取参数字符串

		// 记录函数名信息
		astInfo += fmt.Sprintf("Function Name: %s\r\n", functionName)

		// 解析函数参数
		args := parseArguments(argsString)
		for argName, argValue := range args {
			// 记录每个参数的名称和值
			astInfo += fmt.Sprintf("Argument Name: %s\n", argName)
			astInfo += fmt.Sprintf("Argument Value: %v\n", argValue)
			jsonResult[argName] = argValue // 将参数添加到结果映射中
		}
	}

	return astInfo, jsonResult
}

// parseArguments 从字符串中解析函数参数
// 输入: 参数字符串，如 "first_int={'title': 'First Int'}, second_int={'title': 'Second Int'}"
// 输出: 参数名到参数值的映射
func parseArguments(argsString string) map[string]interface{} {
	result := make(map[string]interface{})

	// 按逗号分割参数，但要注意嵌套结构（如字典、列表等）
	args := splitArgs(argsString)

	// 遍历每个参数
	for _, arg := range args {
		// 按等号分割参数名和参数值，最多分割2次
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			argName := strings.TrimSpace(parts[0])  // 参数名
			argValue := strings.TrimSpace(parts[1]) // 参数值

			// 尝试解析参数值（自动推断类型）
			parsedValue := parseValue(argValue)
			result[argName] = parsedValue
		}
	}

	return result
}

// splitArgs 按逗号分割参数字符串，同时尊重嵌套结构
// 这个函数能够正确处理包含括号和花括号的复杂参数
// 例如: "a=1, b={x:1,y:2}, c=[1,2,3]"
func splitArgs(argsString string) []string {
	var result []string
	var current strings.Builder // 用于构建当前参数
	parenCount := 0             // 括号计数器
	braceCount := 0             // 花括号计数器

	// 逐字符遍历参数字符串
	for _, char := range argsString {
		switch char {
		case '(':
			parenCount++ // 遇到左括号，计数器加1
		case ')':
			parenCount-- // 遇到右括号，计数器减1
		case '{':
			braceCount++ // 遇到左花括号，计数器加1
		case '}':
			braceCount-- // 遇到右花括号，计数器减1
		case ',':
			// 只有在括号和花括号都平衡时，才按逗号分割
			if parenCount == 0 && braceCount == 0 {
				result = append(result, strings.TrimSpace(current.String()))
				current.Reset() // 重置当前参数构建器
				continue
			}
		}
		current.WriteRune(char) // 将当前字符添加到参数中
	}

	// 处理最后一个参数（如果存在）
	if current.Len() > 0 {
		result = append(result, strings.TrimSpace(current.String()))
	}

	return result
}

// parseValue 尝试将字符串值解析为适当的 Go 类型
// 支持的类型: 布尔值、整数、浮点数、JSON 对象、字符串
func parseValue(value string) interface{} {
	// 移除首尾的引号（单引号或双引号）
	value = strings.Trim(value, "'\"")

	// 尝试解析为布尔值
	if value == "true" || value == "false" {
		return value == "true"
	}

	// 尝试解析为整数
	if num, err := strconv.Atoi(value); err == nil {
		return num
	}

	// 尝试解析为浮点数
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num
	}

	// 尝试解析为 JSON 对象（以 { 开头，以 } 结尾）
	if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
		// 先尝试修复 JSON 格式
		repairedValue := repairJSON(value)
		var obj map[string]interface{}
		if err := json.Unmarshal([]byte(repairedValue), &obj); err == nil {
			return obj
		}
	}

	// 如果以上都不匹配，则作为字符串返回
	return value
}

// TryParseJSONObject 清理和格式化 JSON 字符串
// 这个函数尝试修复格式错误的 JSON，包括移除 markdown 代码块、清理格式等
// 返回值: (清理后的 JSON 字符串, 解析后的 JSON 对象)
func TryParseJSONObject(input string) (string, map[string]interface{}) {
	var result map[string]interface{}

	// 首先尝试直接解析 JSON
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return input, result
	}

	// 如果直接解析失败，记录警告信息
	log.Println("Warning: Error decoding faulty json, attempting repair")

	// 使用正则表达式提取 JSON 对象
	// 模式: \{(.*)\} - 匹配花括号及其内容
	pattern := regexp.MustCompile(`\{(.*)\}`)
	matches := pattern.FindStringSubmatch(input)
	if len(matches) >= 2 {
		input = "{" + matches[1] + "}" // 重新构建 JSON 字符串
	}

	// 清理 JSON 字符串（移除常见格式问题）
	input = cleanJSONString(input)

	// 移除 JSON Markdown 代码块框架
	input = removeMarkdownFrame(input)

	// 再次尝试解析清理后的 JSON
	if err := json.Unmarshal([]byte(input), &result); err == nil {
		return input, result
	}

	// 如果仍然失败，尝试修复 JSON
	jsonInfo := repairJSON(input)

	// 尝试解析修复后的 JSON
	if err := json.Unmarshal([]byte(jsonInfo), &result); err != nil {
		log.Printf("Error loading json, json=%s", input)
		return jsonInfo, make(map[string]interface{}) // 返回空映射
	}

	return jsonInfo, result
}

// cleanJSONString 清理常见的 JSON 格式问题
// 处理各种格式错误，如多余的花括号、转义字符等
func cleanJSONString(input string) string {
	// 定义需要替换的字符映射
	replacements := map[string]string{
		"{{":   "{",  // 双花括号替换为单花括号
		"}}":   "}",  // 双花括号替换为单花括号
		"\"[{": "[{", // 修复数组开始标记
		"}]\"": "}]", // 修复数组结束标记
		"\\n":  " ",  // 转义换行符替换为空格
		"\n":   " ",  // 换行符替换为空格
		"\r":   "",   // 回车符删除
		"\\":   " ",  // 反斜杠替换为空格（最后处理）
	}

	// 执行所有替换操作
	for old, new := range replacements {
		input = strings.ReplaceAll(input, old, new)
	}

	// 去除首尾空白字符
	return strings.TrimSpace(input)
}

// removeMarkdownFrame 移除 markdown 代码块
// 处理被 markdown 代码块包围的 JSON 字符串
func removeMarkdownFrame(input string) string {
	input = strings.TrimSpace(input)

	// 移除开头的 ``` 标记
	if strings.HasPrefix(input, "```") {
		input = input[3:]
	}
	// 移除开头的 json 标记（如果存在）
	if strings.HasPrefix(input, "json") {
		input = input[4:]
	}
	// 移除结尾的 ``` 标记
	if strings.HasSuffix(input, "```") {
		input = input[:len(input)-3]
	}

	return strings.TrimSpace(input)
}

// repairJSON 尝试修复格式错误的 JSON
// 这是一个简化版本 - 在实际应用中，你可能想要使用更复杂的 JSON 修复库
// 当前实现包含基本的修复逻辑：移除尾随逗号、修复缺失的引号、转换单引号等
func repairJSON(jsonStr string) string {
	// 基本的 JSON 修复逻辑
	// 这是一个简化实现 - 你可能想要使用适当的 JSON 修复库

	// 移除尾随逗号（在数组或对象结束前）
	// 模式: ,(\s*[}\]]) - 匹配逗号后跟空白字符和 ] 或 }
	re := regexp.MustCompile(`,(\s*[}\]])`)
	jsonStr = re.ReplaceAllString(jsonStr, "$1")

	// 修复键名周围缺失的引号
	// 模式: (\w+): - 匹配字母数字字符后跟冒号
	re = regexp.MustCompile(`(\w+):`)
	jsonStr = re.ReplaceAllString(jsonStr, `"$1":`)

	// 将单引号转换为双引号
	jsonStr = strings.ReplaceAll(jsonStr, "'", "\"")

	return jsonStr
}

// main 函数 - 示例用法
// 演示如何使用这两个主要函数
func main() {
	// 示例 1: 解析函数字符串
	// 展示如何将函数调用字符串转换为 JSON 格式
	functionString := "tool_call(first_int={'title': 'First Int', 'type': 'integer'}, second_int={'title': 'Second Int', 'type': 'integer'})"
	astInfo, jsonResult := Parse(functionString)
	fmt.Println("AST Info:", astInfo)
	fmt.Println("JSON Result:", jsonResult)

	// 示例 2: 解析 JSON 对象
	// 展示如何清理和解析 JSON 字符串
	jsonString := `{"name": "test", "value": 123}`
	cleaned, result := TryParseJSONObject(jsonString)
	fmt.Println("Cleaned:", cleaned)
	fmt.Println("Result:", result)
}
