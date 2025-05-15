package user

import (
	"friedbot/pkg/models"
	"friedbot/pkg/models/schema"
)

func Create(user *schema.User) error {
	return models.DB.Create(user).Error
}

func GetByQQ(qq int) (*schema.User, error) {
	var user *schema.User
	err := models.DB.Where("qq = ?", qq).First(user).Error
	return user, err
}

func UpdateByQQ(user *schema.User) error {
	return models.DB.Where("qq = ?", user.QQ).Updates(user).Error
}

func DeleteByQQ(qq int) error {
	return models.DB.Where("qq = ?", qq).Delete(&schema.User{}).Error
}

func GetAll() ([]*schema.User, error) {
	var users []*schema.User
	err := models.DB.Find(&users).Error
	return users, err
}

func GetByID(id int) (*schema.User, error) {
	var user *schema.User
	err := models.DB.Where("id = ?", id).First(user).Error
	return user, err
}

func UpdateOrCreateByID(user *schema.User) error {
	return models.DB.Save(user).Error
}
