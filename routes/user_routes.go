package routes

import (
	"log"
	"myapp/database"
	"myapp/middleware"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/users", GetUsers)
	app.Get("/users/:email", GetUser)

	// Admin routes (protected with AdminMiddleware)
	app.Post("/users", middleware.AdminMiddleware, CreateUser)
	app.Put("/users/:email/credits", middleware.AdminMiddleware, UpdateCredits)
	app.Put("/users/:email/pause", middleware.AdminMiddleware, PauseUser)
	app.Put("/users/:email/unpause", middleware.AdminMiddleware, UnpauseUser)
	app.Delete("/users/:email", middleware.AdminMiddleware, DeleteUser)
}

// GetUsers fetches all users
func GetUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Find(&users)
	return c.JSON(users)
}

func AddUser(c *fiber.Ctx) error {
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error adding user")
	}

	return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

// CreateUser creates a new user remember this is an admin-only route
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Failed to parse user data: %v", err)
		return c.Status(400).SendString("Bad Request")
	}

	log.Printf("Attempting to create user: %v", user)

	var existingUser models.User
	if err := database.DB.First(&existingUser, "email = ?", user.Email).Error; err == nil {
		log.Printf("User already exists: %v", existingUser)
		return c.Status(400).SendString("User already exists")
	}

	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(500).SendString("Internal Server Error")
	}

	log.Printf("User created successfully: %v", user)
	return c.JSON(user)
}

// this is also an admin only route
func UpdateCredits(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	var creditChange struct {
		Amount int `json:"credits"`
	}
	if err := c.BodyParser(&creditChange); err != nil {
		return c.Status(400).SendString("Bad Request")
	}

	if creditChange.Amount >= 0 {
		user.IncrementCredits(creditChange.Amount)
	} else {
		err := user.DecrementCredits(-creditChange.Amount)
		if err != nil {
			return c.Status(400).SendString(err.Error())
		}
	}

	database.DB.Save(&user)
	return c.JSON(user)
}

// This is also an admin only route
func PauseUser(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	user.PauseAccount()
	database.DB.Save(&user)
	return c.JSON(user)
}

// This is also an admin only route
func UnpauseUser(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	user.UnpauseAccount()
	database.DB.Save(&user)
	return c.JSON(user)
}

// This is also an admin only route
func DeleteUser(c *fiber.Ctx) error {
	email := c.Params("email")
	if err := database.DB.Delete(&models.User{}, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendString("User deleted successfully")
}
