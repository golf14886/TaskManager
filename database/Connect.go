package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	ID       int    `json:"id" gorm:"autoIncrement"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	Status   bool   `json:"status"`
	Lists    []List `json:"lists" gorm:"foreignKey:UserID"`
}

type List struct {
	gorm.Model
	UserID int    `json:"userId"`
	Text   string `json:"text"`
	Cheng  bool   `json:"cheng"`
}

func Connect() (*gorm.DB, error) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		return nil, errors.New("Error loading .env file")
	}

	// Get values from environment variables
	host := envMap["HOST"]
	user := envMap["USER"]
	password := envMap["PASSWORD"]
	dbname := envMap["DBNAME"]
	port := envMap["PORT"]

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("failed to connect database")
	}
	log.Println("Connected to PostgreSQL database")

	err = db.AutoMigrate(&Users{}, &List{})
	if err != nil {
		panic(err)
	}
	return db, nil
}
