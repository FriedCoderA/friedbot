package user

type User struct {
	ID   int    `gorm:"primaryKey"`
	QQ   int    `gorm:"uniqueIndex,not null"`
	Name string `gorm:"not null"`
}

func GetUser(qq int) (*User, error) {
	return &User{
		ID:   1,
		QQ:   123456,
		Name: "test",
	}, nil
}
