package middleware

import (
	"myapp/database"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	email := c.Get("X-User-Email")
	var user models.User

	// Fetch the user from the database using their email
	if err := database.DB.First(&user, "email = ?", email).Error; err != nil {
		return c.Status(404).SendString("User not found")
	}

	if !user.IsAdmin() {
		return c.Status(403).SendString("Access denied: Admins only")
	}

	// Continue to the next handler
	return c.Next()
}
