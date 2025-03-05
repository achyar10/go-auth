package auth

import "github.com/gofiber/fiber/v2"

type AuthController struct {
	Service AuthService
}

func NewAuthController(service AuthService) *AuthController {
	return &AuthController{Service: service}
}

func (ac *AuthController) Register(ctx *fiber.Ctx) error {
	response := ac.Service.Register(ctx)
	return ctx.Status(response.Status).JSON(response)
}

func (ac *AuthController) Login(ctx *fiber.Ctx) error {
	response := ac.Service.Login(ctx)
	return ctx.Status(response.Status).JSON(response)
}

func (ac *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	response := ac.Service.RefreshToken(ctx)
	return ctx.Status(response.Status).JSON(response)
}
