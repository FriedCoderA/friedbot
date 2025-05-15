package models

import (
	"friedbot/pkg/models/schema"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DBPath = "database/test.db"
)

var (
	tables = []interface{}{
		&schema.User{},
	}
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func InitModel() error {
	var err error
	DB, err = gorm.Open(sqlite.Open(DBPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// 模型迁移
	err = DB.AutoMigrate(tables...)
	if err != nil {
		return err
	}

	return nil
}
