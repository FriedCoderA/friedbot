package chat

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"friedbot/pkg/aigc"
	"friedbot/pkg/config"
	"friedbot/pkg/events"
	"friedbot/pkg/kinds"
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
	msgLoadCount      = 50
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
		MaxTokens:        8192,
		FrequencyPenalty: 2,
		PresencePenalty:  2,
		Temperature:      1.2,
	}
	loadCnt := msgLoadCount
	if chat.session.MessageType == kinds.MessageTypePrivate {
		loadCnt += 20
	}
	msgManager := dao.NewMessageManager(chat.session.ID)
	userMessages, err := msgManager.TopN(loadCnt)
	if err != nil {
		slog.Error("ai chat get user messages error", "err", err)
	}
	selfQQ := config.GetBotSettings().QQ
	for _, userMessage := range userMessages {
		if userMessage.CreatedAt.Before(time.Now().Add(-msgExpireDuration)) {
			continue
		}
		if userMessage.UserID == selfQQ {
			req.Messages = append(req.Messages, aigc.NewAssistantMessage(userMessage.Content, "拟人机器人", false, ""))
		} else {
			username := fmt.Sprintf("%s(%d)", userMessage.Sender.GetName(), userMessage.Sender.UserID)
			content := fmt.Sprintf("[%s] %s\n", userMessage.CreatedAt.Format(time.DateTime), userMessage.Content)
			req.Messages = append(req.Messages, aigc.NewUserMessage(username+content, username))
		}
	}
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
		if err != nil {
			slog.Error("chat bot thinking error", "error", err)
			return
		}
		line := strings.Builder{}
		slog.Debug("AI Thinking...")
		reply.Range(func(word string) bool {
			if !strings.Contains(word, "\n") {
				line.WriteString(word)
				return true
			} else {
				if err = onebot.SlowlyReply(chat.session, line.String()); err != nil {
					slog.Error("chat bot reply error", "error", err)
					return false
				}
				line.Reset()
				return true
			}
		})
		if line.Len() > 0 {
			if err = onebot.SlowlyReply(chat.session, line.String()); err != nil {
				slog.Error("chat bot reply error", "error", err)
				return
			}
		}
		chat.state = StateNormal
	}()
	return true, nil
}
