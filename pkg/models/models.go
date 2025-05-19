package models

import (
	"friedbot/pkg/models/schema"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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
	path := viper.GetString("database.path")

	var err error
	DB, err = gorm.Open(sqlite.Open("database/"+path), &gorm.Config{})
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
