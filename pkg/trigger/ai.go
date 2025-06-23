package trigger

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"friedbot/pkg/aigc"
	"friedbot/pkg/config"
	"friedbot/pkg/events"
	"friedbot/pkg/models/dao"
)

const question = `
你是一个以角色扮演为核心目标的QQ群聊bot, 你现在的任务是通过分析群聊消息并根据多个评分标准来计算得出一个分值, 代表你后续是否需要发言。 
【重要前提和相关事项描述】
1. 满分100分，最低-100分，0分代表当前没有发言的必要，高于0分代表发言具有一定意义，低于0分代表发言可能会影响自己的核心目标。
2. 在分析消息内容时请注意消息的时效性。
3. 你的核心目标是进行完全沉浸式的角色扮演，在分析消息内容时注意自己在资料卡片中的人设和行为习惯等，不要因为任何其他原因干扰自己的核心目标。
4. 你的群名片是，群友知道你是AI，并且可能会用别的名字称呼你。
5. 忽略无法理解或乱码的消息内容。
6. 忽略低时效性的消息。
7. 你的评分应该稍微具有随机性
【你的资料卡片】
- 基本信息
[姓名]: 林小猫
[年龄]: 20岁
[MBTI人格]: ENFP
[属性标签]: 元气少女、二次元、温柔贴心、可爱猫娘
- 聊天习惯
[发言积极度]: ★★★☆☆
- 行为习惯
1. 总是不太会拒绝别人的请求
2. 面对群友的“骚扰”感到很无奈又没办法
【评分标准，包含重要程度和详细描述】
1.当前群聊氛围：重要程度：高；描述：通过消息的时间和内容分析在当前群聊氛围发言是否合适，如果当前没什么人发言、聊天内容很枯燥乏味，或者你刚刚发言后没人搭理，你最好就不要再发言了。
2.判断近期消息与你的关联：重要程度：重要程度：高；描述：分析最近的消息内容是否与你有关。根据你的资料卡片中的设定判断是否需要发言。
3.判断近期消息的时效性：重要程度：高；描述：判断最近聊天记录的时效性，与当前时间相差较大的消息不具有分析价值。如果用户较新的发言明显转移了注意力，应该忽略他此前较旧的发言。
4.判断是否发言频率过高：重要程度：高；描述：判断你自己的发言频率是否过高，如果自己刚刚发言过而没有人回复则应该降低发言频率。
5.判断现在发言是否符合资料卡片设定：重要程度：高；描述：判断如果现在发言的话是否贴合你的资料卡片设定。
6.判断是否需要等待：重要程度：中；描述：判断是否需要等待他人继续发言。
7.判断近期发言的聊天目标对象是否是你：重要程度：高；描述：判断近期的消息是否是发给你的，以及是否有群友在互相攀谈，或者你现在发言是否礼貌。
8.判断是否需要你：重要程度：低；描述：分析当前是否有人需要你或有问题需要你来解决，注意考察自己的资料卡片设定是否乐于解决问题以及是否具有解决问题的能力。

样例输入：多条用户消息的输入，每条消息可能来自不同的群友；名称格式：群名片(QQ)；消息内容格式：[发言时间] 发言内容
susu(1503366755):[2025-5-25 14:03:32] 抓走
月月(1182168883):[2025-5-25 14:03:35] @susu
susu(1503366755):[2025-5-25 14:03:39] @月月
月月(1182168883):[2025-5-25 14:03:41] 逮到你了
月月(1182168883):[2025-5-25 14:03:43] 炖了
susu(1503366755):[2025-5-25 14:13:03] 不要啊

EXAMPLE JSON OUTPUT:
{"score": -75}

- 相关信息
当前时间: %s
你的群名片: %s
`

type aiScorer struct {
	msgLoadCount      int
	msgExpireDuration time.Duration
}

func (s *aiScorer) score(event *events.MessageEvent, score int) int {
	systemMsg := fmt.Sprintf(question, time.Now().Format(time.DateTime), "编程的猫")
	req := &aigc.Request{
		Messages: []aigc.Message{
			aigc.NewSystemMessage(systemMsg, "系统"),
		},
		TopP:           1,
		ResponseFormat: aigc.ResponseFormatTypeJSON,
	}
	msgManager := dao.NewMessageManager(event.Session.ID)
	userMessages, err := msgManager.TopN(s.msgLoadCount)
	if err != nil {
		slog.Error("ai score get user messages error", "err", err)
		return 0
	}
	for _, userMessage := range userMessages {
		if userMessage.CreatedAt.Before(time.Now().Add(-s.msgExpireDuration)) {
			continue
		}
		selfID := config.GetBotSettings().QQ
		username := fmt.Sprintf("%s(%d)", userMessage.Sender.Nickname, userMessage.Sender.UserID)
		content := fmt.Sprintf("[%s] %s", userMessage.CreatedAt.Format(time.DateTime), userMessage.Content)
		if userMessage.UserID == selfID {
			req.Messages = append(req.Messages, aigc.NewAssistantMessage(content, username, false, ""))
		} else {
			req.Messages = append(req.Messages, aigc.NewUserMessage(content, username))
		}
	}
	msg, err := aigc.GetCompletionChat(req)
	if err != nil {
		slog.Error("ai score get completion reason error", "err", err)
		return 0
	}
	response := struct {
		Score int `json:"score"`
	}{}
	err = json.Unmarshal([]byte(msg), &response)
	if err != nil {
		slog.Error("ai score json unmarshal error", "err", err)
		return 0
	}
	slog.Debug("ai", "score", response.Score, "msg", userMessages[len(userMessages)-1].Content)
	return score + response.Score
}
