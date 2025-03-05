package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/achyar10/go-auth/src/app/auth"
	"github.com/achyar10/go-auth/src/app/user"
)

// SetupRoutes mengatur semua endpoint dalam aplikasi
func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// Auth
	auth.SetupAuthRoutes(app, db)

	// User
	user.SetupRoutes(app, db)
}
