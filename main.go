package main

import (
	"github.com/achyar10/go-auth/src/app/user"
	"github.com/achyar10/go-auth/src/config"
	"github.com/achyar10/go-auth/src/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Koneksi ke database
	db := config.ConnectDatabase()

	// Inisialisasi Fiber
	app := fiber.New()

	app.Use(helmet.New())

	// Middleware Logging
	app.Use(logger.New())

	// Setup routing
	routes.SetupRoutes(app, db)

	// Jalankan server di port 3000
	db.AutoMigrate(&user.User{})
	app.Listen(":3000")
}
