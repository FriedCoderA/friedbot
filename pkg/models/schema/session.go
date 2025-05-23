package schema

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"
	"strconv"
	"time"

	"friedbot/pkg/config"
	"friedbot/pkg/onebot"
)

type Session struct {
	ID          int64     `gorm:"primaryKey"`
	MessageType string    `json:"message_type" gorm:"not null,index:user,index:group"`
	UserID      int64     `json:"user_id" gorm:"index:user"`
	GroupID     int64     `json:"group_id" gorm:"index:group"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Message struct {
	SessionID   int64     `gorm:"not null"`
	MessageType string    `json:"message_type" gorm:"not null"`
	SubType     string    `json:"sub_type" gorm:"not null"`
	Content     string    `json:"message" gorm:"not null"`
	UserID      int64     `json:"user_id" gorm:"not null"`
	GroupID     int64     `json:"group_id" gorm:""`
	SelfID      int64     `json:"self_id" gorm:"not null"`
	Sender      Sender    `json:"sender" gorm:"type:json"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (m *Message) IsAccess() bool {
	if m.MessageType != onebot.MessageTypePrivate && m.MessageType != onebot.MessageTypeGroup {
		return false
	}
	botSettings := config.GetBotSettings()
	if !slices.Contains(botSettings.GroupWhiteList, strconv.FormatInt(m.GroupID, 10)) {
		return false
	}
	return !slices.Contains(botSettings.UserBlackList, strconv.FormatInt(m.UserID, 10))
}

type Sender struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Card     string `json:"card"`
}

func (s *Sender) Scan(value interface{}) error {
	data, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal Author value")
	}
	var sender Sender
	err := json.Unmarshal(data, &sender)
	if err != nil {
		return err
	}
	*s = sender
	return nil
}

// Value 将Author类型转换为数据库可存储的值
func (s Sender) Value() (driver.Value, error) {
	return json.Marshal(s)
}
