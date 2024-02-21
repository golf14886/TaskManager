package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golf14886/Tasklist/database"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	db, err := database.Connect()
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	app.Post("/signup", func(c *fiber.Ctx) error {
		return Signup(c, db)
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		return SignIn(c, db)
	})


	app.Use(AuthMiddleware)

	// //get list
	app.Get("/:email", func(c *fiber.Ctx) error {
		return GetList(c, db)
	})

	//add list

	app.Post("/:email", func(c *fiber.Ctx) error {
		return AddList(c, db)
	})

	// //EditList
	app.Put("/:email/:id", func(c *fiber.Ctx) error {
		return EditData(c,db)
	})


	app.Delete("/:email/:id", func(c *fiber.Ctx) error {
		return DeleteData(c, db)
	})


	app.Listen(":3000")
}
