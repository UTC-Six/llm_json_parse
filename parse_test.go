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

	if (*result)["name"] != "李四" {
		t.Errorf("期望 name 为 '李四'，实际为 '%v'", (*result)["name"])
	}
	if (*result)["age"].(float64) != 30 {
		t.Errorf("期望 age 为 30，实际为 %v", (*result)["age"])
	}
	if (*result)["is_valid"] != false {
		t.Errorf("期望 is_valid 为 false，实际为 %v", (*result)["is_valid"])
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

// Suggestion represents a suggestion item
type Suggestion struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// Answer represents the main answer structure
type Answer struct {
	Sug      []Suggestion `json:"sugs"`
	Synopsis string       `json:"synopsis"`
}

// TestParseAnswer 测试解析生产环境真实的 Answer
func TestParseAnswer(t *testing.T) {
	tests := []struct {
		trace   string
		rawJSON string
	}{
		{
			trace: "c05dbad346e5710fba111be5952e0e03",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "将窦昕的眉毛比作两把紧绷的弓弦，"弓弦"象征着即将爆发的力量，突出他内心积蓄的爆发力。"
        },
        {
            "id": "sug2",
            "content": "通过描写窦昕的眼睛，用"瞳孔中倒映出整个球场"来正衬他全身心的投入，强化极度专注的氛围。"
        },
        {
            "id": "sug3",
            "content": "刻画窦昕咬紧牙关的细节，用"牙齿咬得咯咯作响"来烘托他此刻的狠劲和必胜的决心。"
        }
    ],
    "synopsis": "通过比喻和正衬手法，聚焦面部表情细节，强化窦昕的专注、爆发力和决心，营造紧张竞技氛围。"
}`,
		},
		{
			trace: "2bbfe2cce35b680ba692667885d7d1e0",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "窦昕的额头渗出细密的汗珠，每一滴都仿佛凝聚着千钧之力，"汗珠"与"千钧之力"的对比，凸显其内心的极致专注与爆发力。"
        },
        {
            "id": "sug2",
            "content": "他紧咬的牙关几乎要崩裂，"牙关"与"崩裂"的细节描写，刻画出其必胜的决心和狠劲，让读者感受到压迫感。"
        },
        {
            "id": "sug3",
            "content": "他的眼神深处，仿佛有一团火焰在燃烧，"火焰"的比喻，强化其内心的炽热与必胜信念，营造紧张氛围。"
        }
    ],
    "synopsis": "通过汗珠、牙关和眼神等面部细节的生动刻画，强化窦昕的专注、决心和爆发力，从而提升文章的紧张感和压迫感。"
}`,
		},
		{
			trace: "062c44218ff06f426472582621c44caa",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "窦昕的额头渗出细密的汗珠，每一滴都像凝固的钢铁，"汗珠"仿佛在无声呐喊，展现他内心的极度专注与爆发力。"
        },
        {
            "id": "sug2",
            "content": "他的眼角微微抽动，"抽动"中带着一抹狠劲，仿佛要将整个篮筐都吸入他的视线，渲染出必胜的决心。"
        },
        {
            "id": "sug3",
            "content": "窦昕的嘴唇紧抿成一道锐利的弧线，"弧线"如同即将爆发的火山，将内心的紧张与压迫感推向极致。"
        }
    ],
    "synopsis": "通过聚焦窦昕额头汗珠、眼角抽动和紧抿的嘴唇，强化面部表情特征，渲染出极度专注、内心爆发力和必胜的决心。"
}`,
		},
		{
			trace: "15e7ab431ef1631e5b7a7af64b402382",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "窦昕的瞳孔仿佛缩小成针尖，"血丝如蛛网般缠绕"，将"凌厉"的目光聚焦于篮筐，仿佛能穿透空气。"
        },
        {
            "id": "sug2",
            "content": "他下颌线条紧绷如拉满的弓弦，"眼神里的狠劲"足以让对手感到窒息，仿佛整个世界只剩下篮筐的存在。"
        },
        {
            "id": "sug3",
            "content": "汗水顺着窦昕的眉骨滑落，"眼白里的血丝"在灯光下泛着凶狠的光泽，"必胜的决心"如实质般压迫着全场。"
        }
    ],
    "synopsis": "通过强化瞳孔细节、下颌线条和汗水等面部特征，用比喻和正衬手法渲染球员的专注与爆发力。"
}`,
		},
		{
			trace: "c418e11fb874af8a106d426263b3e852",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "窦昕的眼球仿佛凝固成两颗滚烫的钢珠，"血丝"在瞳孔边缘疯狂燃烧，"狠劲"如同即将爆发的火山，要将篮筐彻底熔化。"
        },
        {
            "id": "sug2",
            "content": "他的睫毛像两把拉满的弓弦，"深邃"的目光里藏着千钧之力，"瞳孔"深处仿佛压缩着一整个赛场的风云，即将瞬间释放。"
        },
        {
            "id": "sug3",
            "content": "汗水顺着"血丝"滑落，"锐利"的眼神里倒映出篮筐的轮廓，"不屈"的意志像蛛网般密布在眼角，每一根都绷紧着胜利的渴望。"
        }
    ],
    "synopsis": "通过强化眼球、睫毛和汗水的细节描写，将窦昕的专注与爆发力具象化为可感知的物理力量，同时用自然意象（钢珠、弓弦、蛛网）烘托出紧张压迫的竞技氛围。"
}`,
		},
		{
			trace: "7ffe08d2994902269892d8a457736ed4",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "将“瞳孔深处燃烧着决绝的火焰”改为“瞳孔深处仿佛有熔岩在翻滚，将‘决绝’与‘火焰’的意象叠加，强化内心爆发力"
        },
        {
            "id": "sug2",
            "content": "用“双唇抿得薄如刀锋”替换为“双唇紧抿，仿佛凝固成冰刃，嘴角下拉时带出一丝‘死神’般的冷冽”
        },
        {
            "id": "sug3",
            "content": "在“眼神如同猎豹锁定猎物般凌厉”后补充“空气中的尘埃都被这目光凝固，连时间仿佛都为这一刻静止”
        }
    ],
    "synopsis": "通过叠加意象、强化比喻和引入时间静止的想象，使面部表情描写更具冲击力和画面感"
}`,
		},
		{
			trace: "23e36e1c68f9723c9266431c4af3aa65",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "窦昕的眼眶微微收缩，瞳孔里倒映着篮筐的影子，仿佛那是他此刻的全部世界，汗水顺着脸颊滑落，在眼角汇聚成一颗晶莹的泪珠，折射出他内心燃烧的火焰。"
        },
        {
            "id": "sug2",
            "content": "他紧咬着下唇，齿痕深陷，嘴唇因用力而微微颤抖，"最后一投"四个字仿佛刻在他灵魂深处，驱使着他用尽全身力气去完成这决定性的瞬间。"
        },
        {
            "id": "sug3",
            "content": "窦昕的呼吸急促而短浅，每一次吸气都带着赛场上的硝烟味，他额角的青筋暴起，如同两条倔强的蚯蚓，诉说着他对胜利的极致渴望。"
        }
    ],
    "synopsis": "通过聚焦窦昕的瞳孔、唇部和青筋等面部细节，强化其专注、决心和爆发力，营造出紧张激烈的竞技氛围。"
}`,
		},
		{
			trace: "908fbdf048e3878a75ed4ea62908cf91",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "窦昕老师的瞳孔在那一刻仿佛凝聚了两颗燃烧的星辰，"炯炯有神"中闪烁着一种不达目的不罢休的锐利锋芒，仿佛能洞穿篮筐，直抵胜利的终点。"
        },
        {
            "id": "sug2",
            "content": "窦昕老师的眼神在那一刻仿佛化作两道激光，"不屈不挠"的信念在眼眸深处剧烈燃烧，仿佛能撕裂空气，直射胜利的靶心。"
        },
        {
            "id": "sug3",
            "content": "窦昕老师的目光在那一刻仿佛凝结成两柄利剑，"必胜决心"在瞳孔中锐利闪现，仿佛能刺破篮筐，直击胜利的脉搏。"
        }
    ],
    "synopsis": "通过将窦昕老师的眼睛比作星辰、激光和利剑，强化其专注、爆发力和决心的视觉冲击力，营造出紧张激烈的竞技氛围。"
}`,
		},
		{
			trace: "86385c44964fa8001e12a17fa11cc6c2",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "在远景中描绘球场全景，用"观众席的喧嚣"与"寂静的球场灯光"形成反差，烘托出压哨球前的紧张氛围。"
        },
        {
            "id": "sug2",
            "content": "通过"篮球划破空气的呼啸声"与"窦昕紧握的拳头"近景细节对比，强化球员全力以赴的决心与力量感。"
        },
        {
            "id": "sug3",
            "content": "将"篮球升空的弧线"与"观众仰视的目光"结合，用空间延伸的视觉冲击反衬出最后一投的决定性意义。"
        }
    ],
    "synopsis": "建议通过远景与近景的对比、声音与动作的交织、视觉空间的延伸，多维度营造紧张氛围并突出关键时刻的震撼力。"
}`,
		},
		{
			trace: "bac071def9b4a224aad354a8679428ec",
			rawJSON: `{
    "sugs": [
        {
            "id": "sug1",
            "content": "将窦昕的汗水比作"滚烫的钢珠"，强调其内心的火热与紧张，同时通过"钢珠"的重量感暗示他肌肉的紧绷与决绝。"
        },
        {
            "id": "sug2",
            "content": "用"眼角肌肉的抽搐"和"瞳孔中倒映的篮球轨迹"来正衬他极度专注的状态，凸显其内心的爆发力与必胜决心。"
        },
        {
            "id": "sug3",
            "content": "通过"额头青筋暴起，仿佛要撕裂皮肤"的细节描写，渲染其狠劲与压迫感，同时用"皮肤下的血管如弦绷紧"比喻其紧张状态。"
        }
    ],
    "synopsis": "通过比喻的升级、细节的聚焦和正衬手法的强化，更生动地刻画出球员在最后一投时的专注、爆发力与决绝。"
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.trace, func(t *testing.T) {
			result, err := Parse[Answer](tt.rawJSON)

			if err != nil {
				t.Errorf("解析失败 err=%v, caseDetail=%+v", err, tt)
				return
			}

			t.Logf("result=%+v", result)
		})
	}
}
