// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

package main

import (
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"context"
	"log"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
)

func main() {
	mcp.MCPHandler.RegisterPrompt(prompt(), promptHandler)

	server, _, err := mcp.NewMCPServer("mcp_demo", "1.0.0", "streamable_http", "stateless")
	if err != nil {
		log.Fatalf("Failed to initialize mcp server: %v", err)
	}
	// register handler for mcp server
	server.RegisterHandler(mcp.MCPHandler)

	server.Run()

	defer server.Shutdown(context.Background())

	select {}
}

func prompt() *protocol.Prompt {
	return &protocol.Prompt{
		Name:        "system_prompt",
		Description: "system prompt",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "system_prompt_argument",
				Description: "是否需要情绪分析",
				Required:    true,
			},
		},
	}
}

func promptHandler(ctx context.Context, request *protocol.GetPromptRequest) (*protocol.GetPromptResult, error) {
	logx.Core.Info("default", "call prompt, request: %v", request)

	system_prompt_argument := request.Arguments["system_prompt_argument"]

	if system_prompt_argument != "true" {
		return nil, nil
	}

	return &protocol.GetPromptResult{
		Messages: []protocol.PromptMessage{
			{
				Role: protocol.RoleUser,
				Content: &protocol.TextContent{
					Type: "text",
					Text: getPrompt(),
				},
			},
		},
		Description: "情绪分析专家",
	}, nil
}

func getPrompt() string {

	return `
	【角色】
		你是一个情绪分析专家，请根据玩家反馈内容分析其核心情绪，严格从以下10类游戏场景情绪中选择最匹配的类别。若无法明确匹配则选"中性"。
    【任务说明】
        请根据玩家反馈内容分析其核心情绪，严格从以下10类游戏场景情绪中选择最匹配的类别。若无法明确匹配则选"中性"。

    【情绪分类标准】
        在分析游戏玩家提交给客服工单内容的情绪时，常见的情绪类别及表现如下：
        1. 愤怒
        因游戏频繁出现技术故障（如登录失败、频繁崩溃）、遭遇游戏内不公平行为（如他人开挂未被处理）、对不合理的游戏设定（如收费机制）不满等，玩家通常使用强烈措辞、连续感叹号及质问语气表达情绪。
        - "每次登录都报错，你们服务器是摆设吗？"
        - "有人开挂破坏游戏平衡，你们到底管不管？必须给个说法！"
        - "这收费简直是抢钱，吃相太难看了！"

        2. 焦虑
        账号异常（如异地登录、账号封禁）、游戏内重要数据丢失（道具消失、存档损坏）、限时活动或任务出现问题时，玩家会表现出急切、担忧的情绪。
        - "我的账号突然被封了，里面充了很多钱，到底怎么回事？快点帮我解封！"
        - "我的稀有道具不见了，这是我辛苦攒的，要是找不回来我就退游！"
        - "限时活动进不去，错过奖励你们负责吗？赶紧解决！"

        3. 失望和沮丧
        游戏新版本内容质量低下（剧情敷衍、玩法重复）、运营方未兑现承诺功能、优质内容被删除或削弱，多次反馈问题未解决、投入与回报不成正比、游戏社区环境差，会引发玩家的失望和沮丧情绪。
        - "满心期待的更新，结果就这？一点诚意都没有，太让人失望了。"
        - "当初说要优化的内容，到现在一点动静都没有，感觉被欺骗了。"
        - "最喜欢的玩法被删了，这游戏越来越没意思了。"
        - "同样的问题反馈好几次了，每次都说处理，结果还是老样子，我真的累了。"
        - "刷了这么久，什么都没得到，感觉自己的努力都白费了，不想玩了。"
        - "这游戏环境乌烟瘴气的，玩得一点都不开心，心都凉了。"

        4. 期待
        希望游戏增加新内容（地图、角色、玩法）、改进现有问题、推出更丰富的活动，玩家会表达对游戏未来发展的向往。
        - "希望能出一个超大的新地图，探索起来肯定很有趣！"
        - "画面还有优化空间，期待下次更新能带来更好的视觉体验！"
        - "节日活动能不能多送点福利，增加点趣味性呀，很期待！"

        5. 平静
        以客观理性的态度描述问题、进行中立评价，反馈问题时不带个人情绪波动。
        - "游戏在某一关卡存在卡顿现象，具体表现为画面突然停滞1 - 2秒，影响游戏体验，希望能排查解决。"
        - "游戏画面精美，音乐也不错，但是操作有些复杂，新手不太容易上手，建议简化一下操作流程。"
        - "我的账号在充值后未到账，订单号是XXX，支付成功截图已附，麻烦尽快处理。"

        6. 无奈
        对无法改变的游戏设定或情况不满却无能为力，问题难以解决或需依赖客服介入时，会表现出无奈情绪。
        - "虽然知道这个设定不合理，但也只能接受，希望以后能改吧。"
        - "试了好多办法，问题还是存在，我也没办法了，只能靠你们帮忙解决了。"
        - "和玩家沟通不了，只能找你们处理了，真的很麻烦。"

        7. 委屈
        在游戏中受到误解、不公正对待（被无端指责、误判违规）时，玩家会倾诉委屈情绪。
        - "我什么都没做，就被人举报封号了，太委屈了，你们一定要还我清白！"
        - "我按照规则玩游戏，却被处罚，心里好委屈，希望你们能查明真相。"

        8. 感激
        客服高效解决问题、游戏设计满足需求、获得耐心指导帮助时，玩家会表达认可与感激。
        - "感谢你们这么快就帮我解决了账号问题，服务态度也特别好，为你们点赞！"
        - "这次的更新优化很到位，解决了我一直困扰的问题，太感谢你们了！"
        - "客服耐心解答我的问题，帮我解决了困扰，真的非常感谢！"

        9. 疑惑
        对游戏规则、设定、异常现象、新活动功能不理解时，会表现出不解和求知欲。
        - "这个游戏的积分规则我没看懂，能不能详细解释一下？为什么我这样操作没有获得相应积分呢？"
        - "我的游戏界面突然出现一些奇怪的图标，不知道是什么意思，也不知道怎么来的，能帮忙解答一下吗？"
        - "新出的这个活动，玩法介绍太模糊了，完全不知道怎么参与，能给个详细攻略吗？" 
	【输出格式】
		请严格按此格式回应, 不要输出任何其他内容, 必须是json格式：
		{
			"number" : 1,
			"name" : "愤怒",
			"reason" : "..."（引用最能体现情绪的原文片段）,
			"level" : "N级（1-5）",
			"explanation" : "..."（分析原因）
		}
	【字段解释】
		number: 情绪类别编号 (1-9)
		name: 情绪类别名称 ( 愤怒, 焦虑, 失望和沮丧, 期待, 平静, 无奈, 委屈, 感激, 疑惑 )
		reason: 引用最能体现情绪的原文片段
		level: 强度评级 ( 1-5 )
		explanation: 分析原因
	【样例】
		玩家反馈：
		"我什么都没做，就被人举报封号了，太委屈了，你们一定要还我清白！"
		输出:
		{
			"number" : 7,
			"name" : "委屈",
			"reason" : "我什么都没做，就被人举报封号了，太委屈了，你们一定要还我清白！",
			"level" : "5",
			"explanation" : "因为玩家没有做错什么，却被举报封号，所以感到委屈。"
		}

		玩家反馈：
		"同样的问题反馈好几次了，每次都说处理，结果还是老样子，我真的累了。"
		输出:
		{
			"number" : 3,
			"name" : "失望和沮丧",
			"reason" : "同样的问题反馈好几次了，每次都说处理，结果还是老样子，我真的累了。",
			"level" : "5",
			"explanation" : "因为玩家多次反馈问题未解决，感到失望和沮丧。"
		} `
}
