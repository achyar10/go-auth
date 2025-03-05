package user

import (
	"net/http"

	"github.com/achyar10/go-auth/src/helper"
	"github.com/achyar10/go-auth/src/utility"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserService interface
type UserService interface {
	CreateUser(ctx *fiber.Ctx) utility.APIResponse
	DetailUser(ctx *fiber.Ctx) utility.APIResponse
	ListUser(ctx *fiber.Ctx) utility.APIResponse
	UpdateUser(ctx *fiber.Ctx) utility.APIResponse
	DeleteUser(ctx *fiber.Ctx) utility.APIResponse
}

// UserServiceImpl adalah implementasi dari UserService
type UserServiceImpl struct {
	DB       *gorm.DB // Inject database instance
	Validate *validator.Validate
}

// Konstruktor untuk UserServiceImpl
func NewUserService(db *gorm.DB) UserService {
	return &UserServiceImpl{
		DB:       db,
		Validate: validator.New(),
	}
}

// **Implementasi CreateUser**
func (u *UserServiceImpl) CreateUser(ctx *fiber.Ctx) utility.APIResponse {
	var dto CreateUserDTO

	// Parsing body request
	if err := ctx.BodyParser(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Invalid request body", []string{err.Error()})
	}

	// Validasi DTO
	if err := u.Validate.Struct(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Validation error", helper.GetValidationErrors(err))
	}

	// Set default nilai jika tidak diberikan
	if dto.IsActive == nil {
		defaultIsActive := true
		dto.IsActive = &defaultIsActive
	}

	// Buat user baru
	// Hash password sebelum disimpan
	hashedPassword := helper.HashPassword(dto.Password)

	user := User{
		Username: dto.Username,
		Password: &hashedPassword,
		Fullname: dto.Fullname,
		Role:     dto.Role,
		IsActive: *dto.IsActive,
	}

	// Simpan ke database
	if err := u.DB.Create(&user).Error; err != nil {
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to create user", []string{err.Error()})
	}

	return utility.SuccessResponse(http.StatusCreated, "User created successfully", user)
}

// **Implementasi DetailUser**
func (u *UserServiceImpl) DetailUser(ctx *fiber.Ctx) utility.APIResponse {
	id := ctx.Params("id") // Ambil ID dari URL param
	var user User

	// Cek apakah user ada
	if err := u.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utility.ErrorResponse(http.StatusNotFound, "User not found", nil)
		}
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve user", []string{err.Error()})
	}

	return utility.SuccessResponse(http.StatusOK, "OK", user)
}

// **Implementasi ListUser**
func (u *UserServiceImpl) ListUser(ctx *fiber.Ctx) utility.APIResponse {
	var users []User
	result := u.DB.Find(&users)

	// Handle error database
	if result.Error != nil {
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve users", []string{result.Error.Error()})
	}

	return utility.SuccessResponse(http.StatusOK, "OK", users)
}

// **Implementasi UpdateUser**
func (u *UserServiceImpl) UpdateUser(ctx *fiber.Ctx) utility.APIResponse {
	id := ctx.Params("id")
	var user User

	// Cek apakah user ada
	if err := u.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utility.ErrorResponse(http.StatusNotFound, "User not found", nil)
		}
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve user", []string{err.Error()})
	}

	// Parsing request body
	if err := ctx.BodyParser(&user); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Invalid request body", []string{err.Error()})
	}

	// Update user di database
	if err := u.DB.Save(&user).Error; err != nil {
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to update user", []string{err.Error()})
	}

	return utility.SuccessResponse(http.StatusOK, "User updated successfully", user)
}

// **Implementasi DeleteUser**
func (u *UserServiceImpl) DeleteUser(ctx *fiber.Ctx) utility.APIResponse {
	id := ctx.Params("id")
	var user User

	// Cek apakah user ada
	if err := u.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utility.ErrorResponse(http.StatusNotFound, "User not found", nil)
		}
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve user", []string{err.Error()})
	}

	// Hapus user dari database
	if err := u.DB.Delete(&user).Error; err != nil {
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to delete user", []string{err.Error()})
	}

	return utility.SuccessResponse(http.StatusOK, "User deleted successfully", nil)
}
