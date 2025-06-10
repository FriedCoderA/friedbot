package onebot

import "friedbot/pkg/models/schema"

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
