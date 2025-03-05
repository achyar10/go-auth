package middleware

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// BasicAuthMiddleware menangani parsing Basic Auth
func BasicAuthMiddleware(ctx *fiber.Ctx) error {
	// Ambil header Authorization
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Missing Authorization header",
		})
	}

	// Format harus "Basic base64(username:password)"
	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) != 2 || splitHeader[0] != "Basic" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid Authorization header format",
		})
	}

	// Decode base64(username:password)
	decoded, err := base64.StdEncoding.DecodeString(splitHeader[1])
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid Basic Auth encoding",
		})
	}

	// Pisahkan username & password
	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid Basic Auth credentials",
		})
	}

	// Simpan username & password di context agar bisa digunakan di service
	ctx.Locals("username", credentials[0])
	ctx.Locals("password", credentials[1])

	return ctx.Next()
}
