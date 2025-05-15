package userprofile

import (
	"friedbot/pkg/models"
	"friedbot/pkg/models/schema"
)

func UpdateOrCreate(userProfile *schema.UserProfile) error {
	return models.DB.Save(userProfile).Error
}

func GetByUserID(id int) (*schema.UserProfile, error) {
	var userProfile schema.UserProfile
	err := models.DB.Where("user_id = ?", id).First(&userProfile).Error
	return &userProfile, err
}
