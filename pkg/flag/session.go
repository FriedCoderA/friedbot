package flag

import (
	"sync"

	"friedbot/pkg/models/schema"
)

type SessionSync struct {
	data *sync.Map
}

var SessionThinking = &SessionSync{
	data: &sync.Map{},
}

func (s *SessionSync) Pause(session *schema.Session) {
	s.data.Store(session.ID, true)
}

func (s *SessionSync) Resume(session *schema.Session) {
	s.data.Delete(session.ID)
}
