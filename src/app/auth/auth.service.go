package auth

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/achyar10/go-auth/src/app/user"
	"github.com/achyar10/go-auth/src/helper"
	"github.com/achyar10/go-auth/src/utility"
)

// AuthService interface
type AuthService interface {
	Register(ctx *fiber.Ctx) utility.APIResponse
	Login(ctx *fiber.Ctx) utility.APIResponse
	RefreshToken(ctx *fiber.Ctx) utility.APIResponse
}

// AuthServiceImpl struct
type AuthServiceImpl struct {
	DB       *gorm.DB
	Validate *validator.Validate
}

// Konstruktor untuk AuthService
func NewAuthService(db *gorm.DB) AuthService {
	return &AuthServiceImpl{
		DB:       db,
		Validate: validator.New(),
	}
}

// Implementasi Register User
func (a *AuthServiceImpl) Register(ctx *fiber.Ctx) utility.APIResponse {
	var dto RegisterDTO

	// Parsing request body
	if err := ctx.BodyParser(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Invalid request body", []string{err.Error()})
	}

	// Validasi DTO
	if err := a.Validate.Struct(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Validation error", helper.GetValidationErrors(err))
	}

	// Hash password
	hashedPassword := helper.HashPassword(dto.Password)

	// Simpan user baru
	newUser := user.User{
		Username: dto.Username,
		Password: &hashedPassword,
		Fullname: &dto.Fullname,
		Role:     user.Role(dto.Role),
		IsActive: true,
	}

	if err := a.DB.Create(&newUser).Error; err != nil {
		return utility.ErrorResponse(http.StatusInternalServerError, "Failed to create user", []string{err.Error()})
	}

	return utility.SuccessResponse(http.StatusCreated, "User created successfully", newUser)
}

// Implementasi Login User
func (a *AuthServiceImpl) Login(ctx *fiber.Ctx) utility.APIResponse {
	var dto LoginDTO
	var foundUser user.User

	username := ctx.Locals("username").(string)
	password := ctx.Locals("password").(string)

	// Parsing body request
	if err := ctx.BodyParser(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Invalid request body", []string{err.Error()})
	}

	// Validasi DTO
	if err := a.Validate.Struct(&dto); err != nil {
		return utility.ErrorResponse(http.StatusBadRequest, "Validation error", helper.GetValidationErrors(err))
	}

	if username != dto.Username || password != dto.Password {
		return utility.ErrorResponse(http.StatusUnauthorized, "Request body does not match Basic Auth credentials", nil)
	}

	// Cek user di database
	if err := a.DB.Where("username = ?", dto.Username).First(&foundUser).Error; err != nil {
		return utility.ErrorResponse(http.StatusUnauthorized, "username or password wrong", nil)
	}

	// Verifikasi password
	if !helper.CheckPasswordHash(dto.Password, *foundUser.Password) {
		return utility.ErrorResponse(http.StatusUnauthorized, "username or password wrong", nil)
	}

	// Generate JWT token
	token, _ := helper.GenerateJWT(foundUser.Id, foundUser.Username, *foundUser.Fullname, string(foundUser.Role))
	responseData := LoginResponse{
		Id:       foundUser.Id,
		Username: foundUser.Username,
		Fullname: *foundUser.Fullname,
		Role:     string(foundUser.Role),
		Token:    token,
	}
	return utility.SuccessResponse(http.StatusOK, "Login success", responseData)
}

// Implementasi RefreshToken
func (a *AuthServiceImpl) RefreshToken(ctx *fiber.Ctx) utility.APIResponse {
	authHeader := ctx.Get("Authorization")
	if authHeader == "" {
		return utility.ErrorResponse(http.StatusUnauthorized, "Missing Authorization header", nil)
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return utility.ErrorResponse(http.StatusUnauthorized, "Invalid token format", nil)
	}

	tokenString := tokenParts[1]
	token, err := helper.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		return utility.ErrorResponse(http.StatusUnauthorized, "Invalid or expired token", nil)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return utility.ErrorResponse(http.StatusInternalServerError, "Invalid token claims", nil)
	}

	// Generate token baru
	newToken, _ := helper.GenerateJWT(int64(claims["user_id"].(float64)), claims["username"].(string), claims["fullname"].(string), claims["role"].(string))

	return utility.SuccessResponse(http.StatusOK, "Token refreshed", fiber.Map{
		"access_token": newToken,
	})
}
