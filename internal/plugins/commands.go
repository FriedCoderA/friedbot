package plugins

import (
	"fmt"

	"friedbot/pkg/events"
)

type Plugin interface {
	Install() error
	Trigger(event *events.MessageEvent) (bool, error)
	Execute(event *events.MessageEvent) error
}

type EventHandler struct {
	plugins []Plugin
}

func InitPlugins() error {
	handler := &EventHandler{
		plugins: []Plugin{
			&parrotCommand{},
		},
	}
	for _, command := range handler.plugins {
		if err := command.Install(); err != nil {
			return err
		}
	}
	events.Messages.Install(handler)
	return nil
}

func (h *EventHandler) handle(event *events.MessageEvent) (bool, error) {
	for index, command := range h.plugins {
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

func (h *EventHandler) Receive(event *events.MessageEvent) (bool, error) {
	return h.handle(event)
}
