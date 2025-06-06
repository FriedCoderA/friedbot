package events

import (
	"friedbot/pkg/kits/xring"
	"friedbot/pkg/models/schema"
	"friedbot/pkg/xmap"
)

type CommandListener struct {
	queue xmap.XMap[int64, *xring.Ring[*schema.Message]]
}

func NewCommandListener() *CommandListener {
	return &CommandListener{
		queue: xmap.XMap[int64, *xring.Ring[*schema.Message]]{},
	}
}

func (l *CommandListener) Push(session *schema.Session, message *schema.Message) {

}
