package schema

import "time"

type Bot struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int64     `gorm:"index"`
	GroupID   int64     `gorm:"index"`
	Profile   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
