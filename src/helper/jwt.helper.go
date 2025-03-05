package helper

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// Load .env file
func init() {
	_ = godotenv.Load()
}

// GenerateJWT membuat token JWT
func GenerateJWT(userID int64, username string, fullname string, role string) (string, error) {
	// Ambil secret key dari .env
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	// Konversi JWT_EXPIRATION ke int
	jwtExpStr := os.Getenv("JWT_EXPIRATION")
	jwtExp, err := strconv.Atoi(jwtExpStr)

	if err != nil || jwtExp <= 0 {
		jwtExp = 24 // Default ke 24 jam jika tidak ada atau error
	}
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("default_secret") // Fallback jika tidak ada ENV
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"fullname": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * time.Duration(jwtExp)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT memvalidasi token JWT
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	// Ambil secret key dari .env
	jwtSecret := []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		jwtSecret = []byte("default_secret") // Fallback jika tidak ada ENV
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
