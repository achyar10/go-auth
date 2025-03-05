package user

import (
	"github.com/achyar10/go-auth/src/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes mengatur semua routing dalam aplikasi
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Inisialisasi service dan controller
	userService := NewUserService(db)
	userController := NewUserController(userService)

	// Grouping endpoint users
	userRoutes := app.Group("/user")

	// Middleware
	userRoutes.Use(middleware.AuthMiddleware)

	userRoutes.Post("/", userController.CreateUser)
	userRoutes.Get("/", userController.ListUser)
	userRoutes.Get("/:id", userController.DetailUser)
	userRoutes.Put("/:id", userController.UpdateUser)
	userRoutes.Delete("/:id", userController.DeleteUser)
}
