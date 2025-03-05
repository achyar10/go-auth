package routes

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/achyar10/go-auth/src/app/user"
)

// SetupRoutes mengatur semua routing dalam aplikasi
func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Inisialisasi service dan controller
	userService := user.NewUserService(db)
	userController := user.NewUserController(userService)

	// Grouping endpoint users
	userRoutes := app.Group("/user")
	userRoutes.Post("/", userController.CreateUser)
	userRoutes.Get("/", userController.ListUser)
	userRoutes.Get("/:id", userController.DetailUser)
	userRoutes.Put("/:id", userController.UpdateUser)
	userRoutes.Delete("/:id", userController.DeleteUser)
}
