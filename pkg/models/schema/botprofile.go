package schema

import "time"

type BotProfile struct {
	ID        int       `gorm:"primaryKey"`
	Content   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
