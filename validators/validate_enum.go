package validators

import (
	"github.com/dangduoc08/gogo/utils"
	"github.com/go-playground/validator/v10"
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
