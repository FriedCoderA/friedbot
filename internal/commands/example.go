package commands

import (
	"strings"

	"friedbot/pkg/events"
	"friedbot/pkg/onebot"
)

type exampleCommand struct {
	path string
}

func (e *exampleCommand) Install() error {
	e.path = "/example"
	return nil
}

func (e *exampleCommand) Trigger(event *events.MessageEvent) (bool, error) {
	if strings.HasPrefix(event.Message.Content, e.path) {
		return true, nil
	}
	return false, nil
}

func (e *exampleCommand) Execute(event *events.MessageEvent) error {
	msg := onebot.NewReplyMessage(event.Message)
	msg.Content = "hello world"
	return onebot.SendMsg(msg)
}
