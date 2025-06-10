package chat

import (
	"log/slog"

	"friedbot/pkg/events"
	"friedbot/pkg/models/schema"
	"friedbot/pkg/scorer"
	"friedbot/pkg/xmap"
)

const (
	StateNormal int8 = iota
	StateThinking
	StatePaused

	msgLoadCount = 20
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

func (b *chatBot) Thinking(chat *chatSession) (string, error) {
	access := score.Trigger(chat.session)
	if access {
		// msgManager := dao.NewMessageManager(chat.session.ID)
		// userMessages, err := msgManager.TopN(msgLoadCount)
		// if err != nil {
		// 	return "", fmt.Errorf("ai score get user messages error: %v", err)
		// }
		// messages := lo.Map(userMessages, func(item schema.Message, _ int) aigc.Message {
		// 	return aigc.Message(aigc.NewUserMessage(item.Content, item.Sender.GetName()))
		// })
		// reply, err := aigc.GetCompletionChat(&aigc.Request{
		// 	Messages: messages,
		// })
		// if err != nil {
		// 	return "", fmt.Errorf("ai score get reply error: %v", err)
		// }
		reply := ""
		return reply, nil
	}
	chat.state = StateNormal
	return "", nil
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
		_, err := bot.Thinking(chat)
		// reply, err := bot.Thinking(chat)
		if err != nil {
			slog.Error("chat bot thinking error", "error", err)
		}
		// if err = onebot.Reply(chat.session, reply); err != nil {
		// 	slog.Error("chat bot reply error", "error", err)
		// }
		chat.state = StateNormal
	}()
	return true, nil
}
