package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePassword(cb func(validator.FieldError)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value != "" {
			isValidPassword := func(password string) bool {
				hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(password)
				hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
				hasSpecialChar := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(password)
				validLen := len(password) >= 8

				return hasUppercase && hasDigit && hasSpecialChar && validLen
			}

			if !isValidPassword(value) {
				cb(FieldError{
					fieldLevel: fl,
					param:      "",
					val:        value,
				})
			}
		}

		return true
	}
}
