package dao

import (
	"time"

	"friedbot/pkg/models"
	"friedbot/pkg/models/schema"
	"friedbot/pkg/onebot"

	"gorm.io/gorm"
)

type SessionManager struct {
	db *gorm.DB
}

func NewSessionManager() *SessionManager {
	db := models.DB
	return &SessionManager{
		db: db,
	}
}

func (m *SessionManager) Create(session *schema.Session) error {
	return m.db.Create(session).Error
}

func (m *SessionManager) Delete(id int64) error {
	return m.db.Where("id = ?", id).Delete(&schema.Session{}).Error
}

func (m *SessionManager) Get(id int) (*schema.Session, error) {
	var session *schema.Session
	err := m.db.Where("id = ?", id).First(session).Error
	return session, err
}

func (m *SessionManager) GetAll() ([]*schema.Session, error) {
	var sessions []*schema.Session
	err := m.db.Find(&sessions).Error
	return sessions, err
}

func (m *SessionManager) GetOrCreate(msg *schema.Message) (*schema.Session, error) {
	session := &schema.Session{
		MessageType: msg.MessageType,
	}
	tx := m.db.Where("message_type = ?", session.MessageType)
	if msg.MessageType == onebot.MessageTypePrivate {
		session.UserID = msg.UserID
		tx.Where("user_id = ?", session.UserID)
	} else {
		session.GroupID = msg.GroupID
		tx.Where("group_id = ?", session.GroupID)
	}
	err := tx.First(session).Error
	if err == nil {
		return session, nil
	}
	err = m.Create(session)
	return session, err
}

func (m *SessionManager) UpdateOrCreate(session *schema.Session) error {
	return m.db.Save(session).Error
}

type MessageManager struct {
	sessionID int64
	db        *gorm.DB
}

func NewMessageManager(sessionID int64) *MessageManager {
	db := models.DB
	return &MessageManager{
		db:        db,
		sessionID: sessionID,
	}
}

func (m *MessageManager) Create(msg *schema.Message) error {
	msg.SessionID = m.sessionID
	return m.db.Create(msg).Error
}

func (m *MessageManager) TopN(n int) ([]schema.Message, error) {
	var messages []schema.Message
	err := m.db.Where("session_id = ?", m.sessionID).Find(messages).Order("created_at DESC").Limit(n).Error
	return messages, err
}

func (m *MessageManager) Delete(id int64) error {
	return m.db.Where("id = ?", id).Delete(&schema.Session{}).Error
}

func (m *MessageManager) AfterTime(time time.Time) ([]schema.Message, error) {
	var messages []schema.Message
	err := m.db.Where("session_id = ?", m.sessionID).Find(messages).Where("created_at > ?", time).Error
	return messages, err
}
