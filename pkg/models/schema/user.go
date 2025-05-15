package schema

import "time"

type User struct {
	ID        int       `gorm:"primaryKey"`
	QQ        int       `gorm:"uniqueIndex,not null"`
	Name      string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
