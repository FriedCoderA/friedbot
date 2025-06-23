package handlers

import (
	"fmt"

	"friedbot/internal/handlers/parrot"
	"friedbot/pkg/events"
)

type Handler interface {
	Install() error
	Trigger(event *events.MessageEvent) (bool, error)
	Execute(event *events.MessageEvent) error
}

type EventHandler struct {
	handlers []Handler
}

func InitHandlers() error {
	handler := &EventHandler{
		handlers: []Handler{
			&parrot.Parrot{},
		},
	}
	for _, command := range handler.handlers {
		if err := command.Install(); err != nil {
			return err
		}
	}
	events.Messages.Install(handler)
	return nil
}

func (h *EventHandler) handle(event *events.MessageEvent) (bool, error) {
	for index, command := range h.handlers {
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
