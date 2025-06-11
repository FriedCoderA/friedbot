package chat

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"friedbot/pkg/aigc"
	"friedbot/pkg/config"
	"friedbot/pkg/events"
	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"
	"friedbot/pkg/onebot"
	score "friedbot/pkg/scorer"
	"friedbot/pkg/xmap"
)

const (
	StateNormal int8 = iota
	StateThinking
	StatePaused

	msgExpireDuration = time.Hour
	msgLoadCount      = 1
)

var bot *chatBot

type chatSession struct {
	state   int8
	session *schema.Session
}

type chatBot struct {
	sessions *xmap.XMap[int64, *chatSession]
}

func InitChatBot() error {
	bot = &chatBot{
		sessions: xmap.NewXMap[int64, *chatSession](),
	}
	events.Messages.Install(bot)
	return nil
}

func (b *chatBot) Thinking(chat *chatSession) (*aigc.Stream, error) {
	req := &aigc.Request{
		Messages: []aigc.Message{
			aigc.NewSystemMessage(systemPrompt, "聊天提示系统"),
		},
	}
	msgManager := dao.NewMessageManager(chat.session.ID)
	userMessages, err := msgManager.TopN(msgLoadCount)
	if err != nil {
		slog.Error("ai chat get user messages error", "err", err)
	}
	combineMsg := strings.Builder{}
	selfQQ := config.GetBotSettings().QQ
	for _, userMessage := range userMessages {
		if userMessage.CreatedAt.Before(time.Now().Add(-msgExpireDuration)) {
			continue
		}
		if userMessage.UserID == selfQQ {
			if combineMsg.Len() > 0 {
				username := fmt.Sprintf("%s(%d)", userMessage.Sender.GetName(), userMessage.Sender.UserID)
				content := fmt.Sprintf("[%s] %s\n", userMessage.CreatedAt.Format(time.DateTime), userMessage.Content)
				combineMsg.WriteString(username + content)
				req.Messages = append(req.Messages, aigc.NewUserMessage(combineMsg.String(), "群友们"))
				combineMsg.Reset()
			}
			req.Messages = append(req.Messages, aigc.NewAssistantMessage(userMessage.Content, "拟人机器人", false, ""))
		}
		username := fmt.Sprintf("%s(%d)", userMessage.Sender.GetName(), userMessage.Sender.UserID)
		content := fmt.Sprintf("[%s] %s\n", userMessage.CreatedAt.Format(time.DateTime), userMessage.Content)
		combineMsg.WriteString(username + content)
	}
	req.Messages = append(req.Messages, aigc.NewUserMessage(combineMsg.String(), "群友们"))
	reply, err := aigc.GetStreamChat(req)
	if err != nil {
		return nil, fmt.Errorf("ai chat get reply error: %v", err)
	}
	return reply, nil
}

func (b *chatBot) Receive(event *events.MessageEvent) (bool, error) {
	chat, _ := bot.sessions.LoadOrStore(event.Session.ID, &chatSession{
		session: event.Session,
	})
	if chat.state != StateNormal {
		return false, nil
	}
	chat.state = StateThinking
	go func() {
		access := score.Trigger(chat.session)
		if !access {
			chat.state = StateNormal
			return
		}
		reply, err := bot.Thinking(chat)
		chat.state = StateNormal
		if err != nil {
			slog.Error("chat bot thinking error", "error", err)
			return
		}
		line := strings.Builder{}
		reply.Range(func(word string) bool {
			slog.Debug("stream range")
			if !strings.Contains(word, "\n") {
				line.WriteString(word)
				return true
			} else {
				if err = onebot.Reply(chat.session, line.String()); err != nil {
					slog.Error("chat bot reply error", "error", err)
					return false
				}
				line.Reset()
				return true
			}
		})
		if line.Len() > 0 {
			if err = onebot.Reply(chat.session, line.String()); err != nil {
				slog.Error("chat bot reply error", "error", err)
				return
			}
		}
		chat.state = StateNormal
	}()
	return true, nil
}
