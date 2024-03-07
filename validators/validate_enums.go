package validators

import (
	"github.com/dangduoc08/gogo/utils"
	"github.com/go-playground/validator/v10"
)

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
