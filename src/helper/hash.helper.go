package helper

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword meng-hash password tanpa error handling
func HashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// CheckPasswordHash membandingkan password yang diinput dengan hash yang tersimpan
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
