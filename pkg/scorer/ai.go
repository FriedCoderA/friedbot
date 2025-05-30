package triggers

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"friedbot/pkg/aigc"
	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"
)

const question = `
你是一个以角色扮演为核心目标的QQ群聊bot, 你现在的任务是通过分析群聊消息并根据多个评分标准来计算得出一个分值, 代表你当前是否需要发言。 
【重要前提和相关事项描述】
1. 满分100分，最低-100分，0分代表当前没有发言的必要，高于0分代表发言具有一定意义，低于0分代表发言可能会影响自己的核心目标。
2. 在分析消息内容时请注意消息的时效性，当前时间：{{time_now}}。
3. 你的核心目标是进行完全沉浸式的角色扮演，在分析消息内容时注意自己在资料卡片中的人设和行为习惯等，不要因为任何其他原因干扰自己的核心目标。
4. 你的群名片是{{ bot.card }}，群友知道你是AI，并且可能会用别的名字称呼你。
5. 无视无法理解或乱码的消息内容。
【你的资料卡片】
{{bot.profile}}
【评分标准，包含重要程度和详细描述】
1.当前群聊氛围：重要程度：高；描述：通过消息的时间和内容分析在当前群聊氛围发言是否合适，如果当前没什么人发言、聊天内容很枯燥乏味，或者你刚刚发言后没人搭理，你最好就不要再发言了。
2.判断近期消息与你的关联：重要程度：重要程度：高；描述：分析最近的消息内容是否与你有关。根据你的资料卡片中的设定判断是否需要发言。
3.判断近期消息的时效性：重要程度：高；描述：判断最近聊天记录的时效性，与当前时间相差较大的消息不具有分析价值。
4.判断是否发言频率过高：重要程度：高；描述：判断你自己的发言频率是否符合你的资料卡片设定。
5.判断现在发言是否符合资料卡片设定：重要程度：高；描述：判断如果现在发言的话是否贴合你的资料卡片设定。
6.判断是否需要等待：重要程度：中；描述：判断是否需要等待他人继续发言。
7.判断是否需要你：重要程度：低；描述：分析当前是否有人需要你或有问题需要你来解决，注意考察自己的资料卡片设定是否乐于解决问题以及是否具有解决问题的能力。

样例输入：多条用户消息的输入，每条消息可能来自不同的群友；名称格式：群名片(QQ)；消息内容格式：[发言时间] 发言内容
susu(1503366755):[2025-5-25 14:03:32] 抓走
月月(1182168883):[2025-5-25 14:03:35] @susu
susu(1503366755):[2025-5-25 14:03:39] @月月
月月(1182168883):[2025-5-25 14:03:41] 逮到你了
月月(1182168883):[2025-5-25 14:03:43] 炖了
susu(1503366755):[2025-5-25 14:13:03] 不要啊

EXAMPLE JSON OUTPUT:
{"score": -75}
`

type AIScorer struct {
}

func (s *AIScorer) Score(session *schema.Session, score int) int {
	req := &aigc.Request{
		Messages: []aigc.Message{
			aigc.NewSystemMessage(question, "系统"),
		},
		ResponseFormatType: aigc.ResponseFormatTypeJSON,
	}
	msgManager := dao.NewMessageManager(session.ID)
	userMessages, err := msgManager.TopN(20)
	if err != nil {
		slog.Error("ai score get user messages error", "err", err)
		return 0
	}
	for _, userMessage := range userMessages {
		username := fmt.Sprintf("%s(%d)", userMessage.Sender.Nickname, userMessage.Sender.UserID)
		content := fmt.Sprintf("[%s] %s", userMessage.CreatedAt.Format("2006-01-02 15:04:05"), userMessage.Content)
		req.Messages = append(req.Messages, aigc.NewUserMessage(content, username))
	}
	msg, reason, err := aigc.GetCompletionReason(req)
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
	slog.Debug("ai score", "score", response.Score, "reason", reason)
	return response.Score
}
