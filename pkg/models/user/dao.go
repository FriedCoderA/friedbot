package user

func GetUser(qq int) (*User, error) {
	var user User
	err := db.Where("qq = ?", qq).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
