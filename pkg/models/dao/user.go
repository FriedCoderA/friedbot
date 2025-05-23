package dao

import (
	"friedbot/pkg/models"
	"friedbot/pkg/models/schema"

	"gorm.io/gorm"
)

type UserManager struct {
	db *gorm.DB
}

func NewUserManager() *UserManager {
	db := models.DB
	return &UserManager{
		db: db,
	}
}

func (m *UserManager) Create(user *schema.User) error {
	return m.db.Create(user).Error
}

func (m *UserManager) Delete(id int64) error {
	return m.db.Where("id = ?", id).Delete(&schema.User{}).Error
}

func (m *UserManager) Get(id int) (*schema.User, error) {
	var user *schema.User
	err := m.db.Where("id = ?", id).First(user).Error
	return user, err
}

func (m *UserManager) GetAll() ([]*schema.User, error) {
	var users []*schema.User
	err := m.db.Find(&users).Error
	return users, err
}

func (m *UserManager) UpdateOrCreate(user *schema.User) error {
	return m.db.Save(user).Error
}
