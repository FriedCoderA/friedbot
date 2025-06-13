package events

import (
	"log/slog"
	"time"

	"friedbot/pkg/kinds"
	"friedbot/pkg/kits/xring"
	"friedbot/pkg/models/schema"
	"friedbot/pkg/xmap"
)

const MaxMessages = 10

var Messages *MessageEventManager

func InitMessageEvents() error {
	Messages = &MessageEventManager{
		data:     xmap.NewXMap[int64, *xring.Ring[*MessageEvent]](),
		handling: xmap.NewXMap[int64, struct{}](),
	}
	go handle()
	return nil
}

func (e *MessageEventManager) Install(handler Handler) {
	e.installed = append(e.installed, handler)
}

func handle() {
	for {
		time.Sleep(500 * time.Millisecond)
		Messages.data.Range(func(sessionID int64, ring *xring.Ring[*MessageEvent]) bool {
			if Messages.isHandling(sessionID) {
				return false
			}
			Messages.handling.Store(sessionID, struct{}{})
			go func() {
				for {
					event, ok := ring.Pop()
					if !ok {
						break
					}
					for _, handler := range Messages.installed {
						ok, err := handler.Receive(event)
						if err != nil {
							slog.Error("handle message event error", "err", err)
						}
						if ok {
							break
						}
					}
				}
				Messages.handling.Delete(sessionID)
			}()
			return true
		})
	}
}

type Handler interface {
	Receive(*MessageEvent) (bool, error)
}

type MessageEvent struct {
	Session *schema.Session
	Message *schema.Message
}

type MessageEventManager struct {
	data      *xmap.XMap[int64, *xring.Ring[*MessageEvent]]
	handling  *xmap.XMap[int64, struct{}]
	installed []Handler
}

func (e *MessageEventManager) Push(session *schema.Session, message *schema.Message) {
	ring, _ := e.data.LoadOrStore(session.ID, xring.NewRing[*MessageEvent](MaxMessages))
	if session.State == kinds.SessionStatePause {
		return
	}
	ring.Push(&MessageEvent{
		Session: session,
		Message: message,
	})
}

func (e *MessageEventManager) isHandling(sessionID int64) bool {
	_, ok := e.handling.Load(sessionID)
	return ok
}
