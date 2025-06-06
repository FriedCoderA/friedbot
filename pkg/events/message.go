package events

import "friedbot/pkg/models/schema"

const MaxMessages = 10000

var MessageEvents = &MessageChannel{
	data: make(chan *MessageEvent, MaxMessages),
}

type MessageEvent struct {
	Session *schema.Session
	Message *schema.Message
}

type MessageChannel struct {
	data      chan *MessageEvent
	installed []Listener[*MessageEvent]
}

func (e *MessageChannel) Push(session *schema.Session, message *schema.Message) {
	e.data <- &MessageEvent{
		Session: session,
		Message: message,
	}
}
