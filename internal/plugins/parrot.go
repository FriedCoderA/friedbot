package plugins

import (
	"friedbot/pkg/config"
	"friedbot/pkg/events"
	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"
	"friedbot/pkg/onebot"

	"github.com/samber/lo"
	"github.com/samber/lo/mutable"
)

type parrotCommand struct {
}

func (e *parrotCommand) Install() error {
	return nil
}

func (e *parrotCommand) Trigger(event *events.MessageEvent) (bool, error) {
	msgManager := dao.NewMessageManager(event.Session.ID)
	messages, err := msgManager.TopN(15)
	if err != nil {
		return false, err
	}
	if len(messages) < 3 {
		return false, nil
	}
	mutable.Reverse(messages)
	parrot := messages[0].Content
	for _, message := range messages[:3] {
		if message.Content != parrot {
			return false, nil
		}
	}
	if len(messages) > 3 {
		selfID := config.GetBotSettings().QQ
		if lo.ContainsBy(messages, func(item schema.Message) bool {
			return item.UserID == selfID && item.Content == parrot
		}) {
			return false, nil
		}
	}
	return true, nil
}

func (e *parrotCommand) Execute(event *events.MessageEvent) error {
	return onebot.Reply(event.Session, event.Message.Content)
}
