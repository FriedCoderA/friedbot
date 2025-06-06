package events

type Listener[Event any] interface {
	Receive(Event) bool
}
