package database

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *Users) error {

	costStr := os.Getenv("BCRYPT_COST")

	cost, err := strconv.Atoi(costStr)
	if err != nil {
		return errors.New("Failed to convert BCRYPT_COST to integer")
	}
	if cost <= 0 {
		return errors.New("BCRYPT_COST must be a positive integer")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), cost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	if err := db.Create(&user).Error; err != nil {
		fmt.Println("Failed to create user:", err)
		return err
	}
	return nil
}
