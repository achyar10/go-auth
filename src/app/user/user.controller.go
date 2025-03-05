package user

import (
	"github.com/gofiber/fiber/v2"
)

// UserController struct
type UserController struct {
	Service UserService
}

// NewUserController adalah constructor untuk UserController
func NewUserController(service UserService) *UserController {
	return &UserController{Service: service}
}

// CreateUser menangani pembuatan pengguna baru
func (uc *UserController) CreateUser(ctx *fiber.Ctx) error {
	response := uc.Service.Create(ctx)
	return ctx.Status(response.Status).JSON(response)
}

// DetailUser menangani pengambilan detail pengguna berdasarkan ID
func (uc *UserController) DetailUser(ctx *fiber.Ctx) error {
	response := uc.Service.Detail(ctx)
	return ctx.Status(response.Status).JSON(response)
}

// ListUser menangani pengambilan daftar semua pengguna
func (uc *UserController) ListUser(ctx *fiber.Ctx) error {
	response := uc.Service.List(ctx)
	return ctx.Status(response.Status).JSON(response)
}

// UpdateUser menangani pembaruan informasi pengguna berdasarkan ID
func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {
	response := uc.Service.Update(ctx)
	return ctx.Status(response.Status).JSON(response)
}

// DeleteUser menangani penghapusan pengguna berdasarkan ID
func (uc *UserController) DeleteUser(ctx *fiber.Ctx) error {
	response := uc.Service.Delete(ctx)
	return ctx.Status(response.Status).JSON(response)
}
