package commands

import (
	"friedbot/pkg/models/schema"
)

var handler *Handler

type Command interface {
	Install() error
	Trigger(session *schema.Session) (bool, error)
	Execute(session *schema.Session) error
}

type Handler struct {
	commands []Command
}

func InitCommands() {
	handler = &Handler{
		commands: []Command{},
	}
	for _, command := range handler.commands {
		if err := command.Install(); err != nil {
			panic(err)
		}
	}
}
