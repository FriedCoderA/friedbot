package commands

import (
	"fmt"
	"strings"

	"friedbot/pkg/events"
)

type Command interface {
	Install() error
	Trigger(event *events.MessageEvent) (bool, error)
	Execute(event *events.MessageEvent) error
}

type Handler struct {
	commands []Command
}

func InitCommands() error {
	handler := &Handler{
		commands: []Command{
			&exampleCommand{},
		},
	}
	for _, command := range handler.commands {
		if err := command.Install(); err != nil {
			return err
		}
	}
	events.Messages.Install(handler)
	return nil
}

func (h *Handler) handle(event *events.MessageEvent) (bool, error) {
	for index, command := range h.commands {
		ok, err := command.Trigger(event)
		if err != nil {
			return false, fmt.Errorf("error triggering command: %v, command_index: %d,session_id:%d, msg: %s",
				err, index, event.Session.ID, event.Message.Content)
		}
		if ok {
			err = command.Execute(event)
			if err != nil {
				return true, fmt.Errorf("error executing command: %v, command_index: %d,session_id:%d, msg: %s",
					err, index, event.Session.ID, event.Message.Content)
			}
		}
	}
	return false, nil
}

func (h *Handler) Receive(event *events.MessageEvent) (bool, error) {
	if len(event.Message.Content) < 2 || !strings.HasPrefix(event.Message.Content, "/") {
		return false, nil
	}
	return h.handle(event)
}
