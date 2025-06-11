package onebot

import (
	"time"

	"friedbot/pkg/config"
	"friedbot/pkg/models/dao"
	"friedbot/pkg/models/schema"
)

type Message struct {
	UserID     int64  `json:"user_id"`
	GroupID    int64  `json:"group_id"`
	Content    string `json:"message"`
	AutoEscape bool   `json:"auto_escape"`
}

func NewReplyMessage(msg *schema.Message) *Message {
	return &Message{
		UserID:  msg.UserID,
		GroupID: msg.GroupID,
	}
}

func SendPrivateMsg(message *Message) error {
	req := &request{
		path: "/send_private_msg",
		body: message,
	}
	_, err := req.Post()
	return err

}

func SendGroupMsg(message *Message) error {
	req := &request{
		path: "/send_group_msg",
		body: message,
	}
	_, err := req.Post()
	return err
}

func SendMsg(message *Message) error {
	if message.GroupID != 0 {
		return SendGroupMsg(message)
	} else {
		return SendPrivateMsg(message)
	}
}

func Reply(session *schema.Session, message string) error {
	if message == "" {
		return nil
	}
	time.Sleep(time.Millisecond * time.Duration(200*len(message)))
	selfID := config.GetBotSettings().QQ
	msgManager := dao.NewMessageManager(session.ID)
	err := msgManager.Create(&schema.Message{
		Content: message,
		UserID:  selfID,
	})
	if err != nil {
		return err
	}
	return SendMsg(&Message{
		Content: message,
		GroupID: session.GroupID,
		UserID:  session.UserID,
	})
}
