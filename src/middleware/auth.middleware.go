package middleware

import (
	"strings"

	"github.com/achyar10/go-auth/src/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware untuk melindungi route dengan JWT
func AuthMiddleware(ctx *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Missing token",
		})
	}

	// Format token harus "Bearer {token}"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid token format",
		})
	}

	tokenString := tokenParts[1]

	// Validasi token JWT
	token, err := helper.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid or expired token",
		})
	}

	// Extract claims dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid token claims",
		})
	}

	// Simpan data user ke context
	ctx.Locals("user_id", claims["user_id"])
	ctx.Locals("username", claims["username"])
	ctx.Locals("fullname", claims["fullname"])
	ctx.Locals("role", claims["role"])

	return ctx.Next()
}
