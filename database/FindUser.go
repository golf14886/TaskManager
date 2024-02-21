package database

import (
	"gorm.io/gorm"
)

func FindUser(db *gorm.DB, Email string) (*Users, error) {
	user := &Users{}
	result := db.Where("email = ?", Email).First(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
