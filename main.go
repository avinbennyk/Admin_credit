package main

import (
	"myapp/database"
	"myapp/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Connect to the database
	database.ConnectDB()

	// Create a new Fiber application
	app := fiber.New()

	// Setup routes
	routes.SetupRoutes(app)

	// Start the server
	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
