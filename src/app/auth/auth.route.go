package auth

import (
	"github.com/achyar10/go-auth/src/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupAuthRoutes mengatur routing untuk authentication
func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	authService := NewAuthService(db)
	authController := NewAuthController(authService)

	authRoutes := app.Group("/auth")
	authRoutes.Post("/register", authController.Register)
	authRoutes.Post("/login", middleware.BasicAuthMiddleware, authController.Login)
	authRoutes.Get("/refresh", middleware.AuthMiddleware, authController.RefreshToken)
}
