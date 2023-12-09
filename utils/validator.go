package utils

import (
	"fmt"

	"github.com/dangduoc08/gooh/utils"
	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func ValidateEnum(enums []string, cb func(error)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if !utils.ArrIncludes[string](enums, value) {
			cb(fmt.Errorf("Field: %v, Error: must be in %v", fl.FieldName(), enums))
		}

		return true
	}
}

func ValidatePhone(countryCode string, cb func(error)) validator.Func {
	return func(fl validator.FieldLevel) bool {
		num, _ := phonenumbers.Parse(fl.Field().String(), countryCode)
		if !phonenumbers.IsValidNumber(num) {
			cb(fmt.Errorf("Field: %v, Error: must be valid", fl.FieldName()))
		}

		return true
	}
}
