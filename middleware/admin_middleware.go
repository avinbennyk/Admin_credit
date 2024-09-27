package middleware

import (
	"myapp/database"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

// AdminMiddleware ensures that only admins can modify users
func AdminMiddleware(c *fiber.Ctx) error {
	email := c.Get("X-User-Email") // Replace with actual auth logic
	var user models.User

	// Fetch the user from the database using their email
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	// Check if the user is an admin
	if !user.IsAdmin() {
		return c.Status(403).SendString("Access denied: Admins only")
	}

	// Continue to the next middleware/handler if the user is an admin
	return c.Next()
}
