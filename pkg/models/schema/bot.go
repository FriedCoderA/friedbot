package schema

import "time"

type Bot struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int64     `gorm:"index"`
	GroupID   int64     `gorm:"index"`
	Profile   int64     `gorm:"not null, default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
