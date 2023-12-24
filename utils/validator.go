package utils

import (
	"regexp"

	"github.com/dangduoc08/gooh/utils"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func ValidateEnum(enums []string, cb func(validator.FieldError)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value != "" && !utils.ArrIncludes[string](enums, value) {
			cb(FieldError{
				fieldLevel: fl,
				param:      enums,
				val:        value,
			})
		}

		return true
	}
}

func ValidateEnums(enums []string, cb func(validator.FieldError)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		totalValue := fl.Field().Len()
		values := fl.Field().Slice(0, totalValue).Interface().([]string)
		for _, value := range values {
			if value != "" && !utils.ArrIncludes[string](enums, value) {
				cb(FieldError{
					fieldLevel: fl,
					param:      enums,
					val:        value,
				})
			}
		}

		return true
	}
}

func ValidatePhone(countryCode string, cb func(validator.FieldError)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value != "" {
			num, _ := phonenumbers.Parse(value, countryCode)
			if !phonenumbers.IsValidNumber(num) {
				cb(FieldError{
					fieldLevel: fl,
					param:      countryCode,
					val:        value,
				})
			}
		}

		return true
	}
}

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
