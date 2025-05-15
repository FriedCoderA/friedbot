package schema

import "time"

type UserProfile struct {
	ID        int       `gorm:"primaryKey"`
	UserID    int       `gorm:"uniqueIndex,not null"`
	Content   string    `gorm:"not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
