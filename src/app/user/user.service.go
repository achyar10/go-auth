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
	Create(ctx *fiber.Ctx) utility.APIResponse
	Detail(ctx *fiber.Ctx) utility.APIResponse
	List(ctx *fiber.Ctx) utility.APIResponse
	Update(ctx *fiber.Ctx) utility.APIResponse
	Delete(ctx *fiber.Ctx) utility.APIResponse
	ResetPassword(ctx *fiber.Ctx) utility.APIResponse
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

// Implementasi ListUser
func (u *UserServiceImpl) List(ctx *fiber.Ctx) utility.APIResponse {
	var users []User

	// Gunakan helper untuk query params
	query := helper.ParseQueryParams(ctx)

	// Tentukan field yang bisa dicari dalam tabel `users`
	searchableFields := []string{"username", "fullname", "role"}

	// Gunakan helper ApplyFiltersAndPagination
	paginatedResult := helper.ApplyFiltersAndPagination(u.DB, &users, query, searchableFields)

	// Generate metadata
	metadata := helper.GenerateMetadata(query, paginatedResult.TotalCount, paginatedResult.PageCount)

	// Response dengan metadata
	responseData := map[string]interface{}{
		"records":  paginatedResult.Records,
		"metadata": metadata,
	}

	return utility.SuccessResponse(http.StatusOK, "OK", responseData)
}

// Implementasi CreateUser
func (u *UserServiceImpl) Create(ctx *fiber.Ctx) utility.APIResponse {
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

// Implementasi DetailUser
func (u *UserServiceImpl) Detail(ctx *fiber.Ctx) utility.APIResponse {
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

// Implementasi UpdateUser
func (u *UserServiceImpl) Update(ctx *fiber.Ctx) utility.APIResponse {
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

// Implementasi DeleteUser
func (u *UserServiceImpl) Delete(ctx *fiber.Ctx) utility.APIResponse {
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

func (u *UserServiceImpl) ResetPassword(ctx *fiber.Ctx) utility.APIResponse {
	id := ctx.Params("id")
	var dto ResetPasswordUserDTO

	// Validasi input DTO
	if err := u.Validate.Struct(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Validation error", helper.GetValidationErrors(err))
	}

	// Cek apakah user ada
	var user User
	if err := u.DB.Select("id").First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utility.ErrorResponse(http.StatusNotFound, "User not found", nil)
		}
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to retrieve user", []string{err.Error()})
	}

	// Update password dengan hashing langsung di query
	if err := u.DB.Model(&User{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"password": helper.HashPassword(dto.NewPassword),
		}).Error; err != nil {
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to reset password", []string{err.Error()})
	}

	return utility.SuccessResponse(http.StatusOK, "Password reset successfully", nil)
}
