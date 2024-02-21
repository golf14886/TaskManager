package database

import "gorm.io/gorm"

func GetListByID(db *gorm.DB, id int) ([]List, error) {
	lists := []List{}
	err := db.Where("user_id =? ", id).Find(&lists).Error
	if err != nil {
		return nil, err
	}
	return lists, nil
}
