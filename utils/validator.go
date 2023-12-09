package utils

import (
	"github.com/dangduoc08/gooh/utils"
	"github.com/go-playground/validator/v10"
)

func ValidateEnum(enums ...string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		return utils.ArrIncludes[string](enums, value)
	}
}
