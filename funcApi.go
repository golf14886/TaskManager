package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/golf14886/Tasklist/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Signup(c *fiber.Ctx, db *gorm.DB) error {

	newUser := &database.Users{}
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}
	if err := database.CreateUser(db, newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
	})
}

func SignIn(c *fiber.Ctx, db *gorm.DB) error {
	user := database.Users{}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	foundUser, err := database.FindUser(db, user.Email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Incorrect password",
		})
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_email": foundUser.Email,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := []byte(os.Getenv("MY_SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to generate token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:  "jwt",
		Value: tokenString,
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User login successful",
		"token":   tokenString,
	})
}

func GetList(c *fiber.Ctx, db *gorm.DB) error {

	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing email parameter")
	}
	user, err := database.FindUser(db, email)
	if err != nil {
		return err
	}

	lists, err := database.GetListByID(db, user.ID)

	if err != nil {
		log.Println("Failed to fetch lists:", err)
		return err
	}

	fmt.Println("Lists:", lists)

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"user":  user,
		"lists": lists})
}

func AddList(c *fiber.Ctx, db *gorm.DB) error {
	email := c.Params("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing email parameter")
	}

	// ค้นหาผู้ใช้โดยใช้อีเมล
	user, err := database.FindUser(db, email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// อ่านข้อมูลจาก body ของคำขอ
	newlist := &database.List{}

	if err := c.BodyParser(&newlist); err != nil {
		return err
	}

	// เพิ่มรายการใหม่ลงในฐานข้อมูล
	// ให้รายการใหม่นี้เชื่อมกับผู้ใช้ที่พบ
	newlist.UserID = user.ID

	if err := db.Create(&newlist).Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "List added successfully",
		"list":    newlist,
	})
}

func EditData(c *fiber.Ctx, db *gorm.DB) error {

	newList := database.List{}

	//url information
	email := c.Params("email")
	idStr := c.Params("id")

	if err := c.BodyParser(&newList); err != nil {
		return err
	}

	//User information
	user, err := database.FindUser(db, email)
	if err != nil {
		return err
	}

	//all list for user
	lists, err := database.GetListByID(db, user.ID)
	if err != nil {
		return err
	}

	// loop listID == idStr

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	for i, list := range lists {
		if list.ID == uint(id) {

			lists[i].Text = newList.Text
			lists[i].Cheng = newList.Cheng

			if err := db.Save(&lists[i]).Error; err != nil {
				return err
			}
			break
		}
	}

	return c.JSON(fiber.Map{
		"update":   "successful",
		"lists":    lists,
		"id":       idStr,
		"newLists": newList,
	})

}

func DeleteData(c *fiber.Ctx, db *gorm.DB) error {
	// ดึง email และ id จาก URL
	email := c.Params("email")
	idStr := c.Params("id")

	// แปลง id จาก string เป็นตัวเลข
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return err
	}

	// ค้นหาข้อมูลผู้ใช้จากฐานข้อมูลโดยใช้อีเมล
	user, err := database.FindUser(db, email)
	if err != nil {
		return err
	}

	// ลบข้อมูลที่มี id ที่ระบุออกจากฐานข้อมูล
	if err := db.Where("id = ? AND user_id = ?", id, user.ID).Delete(database.List{}).Error; err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Delete successful",
	})
}
