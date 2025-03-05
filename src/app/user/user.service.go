package user

import (
	"net/http"

	"github.com/achyar10/go-auth/src/helper"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// UserService interface
type UserService interface {
	CreateUser(ctx *fiber.Ctx) map[string]interface{}
	DetailUser(ctx *fiber.Ctx) map[string]interface{}
	ListUser(ctx *fiber.Ctx) map[string]interface{}
	UpdateUser(ctx *fiber.Ctx) map[string]interface{}
	DeleteUser(ctx *fiber.Ctx) map[string]interface{}
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
func (u *UserServiceImpl) CreateUser(ctx *fiber.Ctx) map[string]interface{} {
	var dto CreateUserDTO

	// Parsing body request
	if err := ctx.BodyParser(&dto); err != nil {
		return fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"error":   []string{err.Error()},
		}
	}

	// Validasi DTO
	if err := u.Validate.Struct(&dto); err != nil {
		return fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Validation error",
			"error":   helper.GetValidationErrors(err),
		}
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
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to create user",
			"error":   []string{err.Error()},
		}
	}

	return fiber.Map{
		"status":  http.StatusCreated,
		"message": "User created successfully",
		"data":    user,
	}
}

// **Implementasi DetailUser**
func (u *UserServiceImpl) DetailUser(ctx *fiber.Ctx) map[string]interface{} {
	id := ctx.Params("id") // Ambil ID dari URL param
	var user User

	// Cek apakah user ada
	if err := u.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.Map{
				"status":  http.StatusNotFound,
				"message": "Data not found",
			}
		}
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Internal server error",
			"error":   []string{err.Error()},
		}
	}

	return fiber.Map{
		"status":  http.StatusOK,
		"message": http.StatusText(http.StatusOK),
		"data":    user,
	}
}

// **Implementasi ListUser**
func (u *UserServiceImpl) ListUser(ctx *fiber.Ctx) map[string]interface{} {
	var users []User
	result := u.DB.Find(&users)

	// Handle error database
	if result.Error != nil {
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to retrieve users",
			"error":   result.Error.Error(),
		}
	}

	return fiber.Map{
		"status":  http.StatusOK,
		"message": http.StatusText(http.StatusOK),
		"data":    users,
	}
}

// **Implementasi UpdateUser**
func (u *UserServiceImpl) UpdateUser(ctx *fiber.Ctx) map[string]interface{} {
	id := ctx.Params("id")
	var user User

	// Cek apakah user ada
	if err := u.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.Map{
				"status":  http.StatusNotFound,
				"message": "Data not found",
			}
		}
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to retrieve user",
			"error":   []string{err.Error()},
		}
	}

	// Parsing request body
	if err := ctx.BodyParser(&user); err != nil {
		return fiber.Map{
			"status":  http.StatusBadRequest,
			"message": "Invalid request body",
			"error":   []string{err.Error()},
		}
	}

	// Update user di database
	if err := u.DB.Save(&user).Error; err != nil {
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to update user",
			"error":   []string{err.Error()},
		}
	}

	return fiber.Map{
		"status":  http.StatusOK,
		"message": "User updated successfully",
		"data":    user,
	}
}

// **Implementasi DeleteUser**
func (u *UserServiceImpl) DeleteUser(ctx *fiber.Ctx) map[string]interface{} {
	id := ctx.Params("id")
	var user User

	// Cek apakah user ada
	if err := u.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.Map{
				"status":  http.StatusNotFound,
				"message": "User not found",
			}
		}
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to retrieve user",
			"error":   []string{err.Error()},
		}
	}

	// Hapus user dari database
	if err := u.DB.Delete(&user).Error; err != nil {
		return fiber.Map{
			"status":  http.StatusInternalServerError,
			"message": "Failed to delete user",
			"error":   []string{err.Error()},
		}
	}

	return fiber.Map{
		"status":  http.StatusOK,
		"message": "User deleted successfully",
	}
}
