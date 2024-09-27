package routes

import (
	"log"
	"myapp/database"
	"myapp/middleware"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes defines the user-related routes
func SetupRoutes(app *fiber.App) {
	// Public routes (e.g., fetching user data)
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

	// Parse the request body into the User model
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad Request")
	}

	// Insert the user into the database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error adding user")
	}

	return c.JSON(user)
}

// GetUser fetches a single user by email
func GetUser(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.JSON(user)
}

// CreateUser creates a new user (admin-only)
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		log.Printf("Failed to parse user data: %v", err)
		return c.Status(400).SendString("Bad Request")
	}

	log.Printf("Attempting to create user: %v", user)

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.First(&existingUser, "email = ?", user.Email).Error; err == nil {
		log.Printf("User already exists: %v", existingUser)
		return c.Status(400).SendString("User already exists")
	}

	// Create the new user
	if err := database.DB.Create(&user).Error; err != nil {
		log.Printf("Error creating user: %v", err)
		return c.Status(500).SendString("Internal Server Error")
	}

	log.Printf("User created successfully: %v", user)
	return c.JSON(user)
}

// UpdateCredits updates user credits (admin-only)
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

// PauseUser pauses the user account (admin-only)
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

// UnpauseUser reactivates the user account
func UnpauseUser(c *fiber.Ctx) error {
	email := c.Params("email")
	var user models.User
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	user.UnpauseAccount()   // Call the method to unpause the account
	database.DB.Save(&user) // Save the updated user back to the database
	return c.JSON(user)     // Return the updated user
}

// DeleteUser deletes a user (admin-only)
func DeleteUser(c *fiber.Ctx) error {
	email := c.Params("email")
	if err := database.DB.Delete(&models.User{}, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}
	return c.SendString("User deleted successfully")
}
