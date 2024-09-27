package main

import (
	"myapp/database"
	"myapp/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectDB()

	app := fiber.New()

	routes.SetupRoutes(app)

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
