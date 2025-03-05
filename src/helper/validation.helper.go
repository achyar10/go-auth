package helper

import (
	"github.com/go-playground/validator/v10"
)

// GetValidationErrors mengembalikan error validasi dalam bentuk array
func GetValidationErrors(err error) []string {
	var errors []string
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrs {
			var errorMessage string

			switch fieldErr.Tag() {
			case "required":
				errorMessage = fieldErr.Field() + " harus diisi"
			case "min":
				errorMessage = fieldErr.Field() + " minimal harus " + fieldErr.Param() + " karakter"
			case "max":
				errorMessage = fieldErr.Field() + " maksimal " + fieldErr.Param() + " karakter"
			case "oneof":
				errorMessage = fieldErr.Field() + " harus salah satu dari " + fieldErr.Param()
			default:
				errorMessage = fieldErr.Field() + " tidak valid"
			}

			errors = append(errors, errorMessage)
		}
	}
	return errors
}
