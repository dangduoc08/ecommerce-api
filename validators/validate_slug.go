package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateSlug(cb func(validator.FieldError)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()

		if value != "" {
			pattern := `^[a-z0-9]+(?:-[a-z0-9]+)*$`
			match, _ := regexp.MatchString(pattern, value)
			if !match {
				cb(FieldError{
					fieldLevel: fl,
					val:        value,
				})
			}
		}

		return true
	}
}
