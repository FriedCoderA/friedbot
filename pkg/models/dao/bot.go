package dao

import (
	"friedbot/pkg/models"
	"friedbot/pkg/models/schema"

	"gorm.io/gorm"
)

type BotManager struct {
	db *gorm.DB
}

func NewBotManager() *BotManager {
	db := models.DB
	return &BotManager{
		db: db,
	}
}

func (m *BotManager) Create(bot *schema.Bot) error {
	return m.db.Create(bot).Error
}

func (m *BotManager) Delete(id int64) error {
	return m.db.Where("id = ?", id).Delete(&schema.Bot{}).Error
}

func (m *BotManager) Updates(bot *schema.User) error {
	return m.db.Updates(bot).Error
}

func (m *BotManager) Get(id int) (*schema.Bot, error) {
	var bot *schema.Bot
	err := m.db.Where("id = ?", id).First(bot).Error
	return bot, err
}

func (m *BotManager) GetAll() ([]*schema.Bot, error) {
	var bots []*schema.Bot
	err := m.db.Find(&bots).Error
	return bots, err
}

func (m *BotManager) UpdateOrCreate(bot *schema.Bot) error {
	return m.db.Save(bot).Error
}
