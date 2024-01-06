package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

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
