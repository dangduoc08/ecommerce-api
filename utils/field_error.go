package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	validator.FieldError
	fieldLevel validator.FieldLevel
	param      any
	val        any
}

func (fe FieldError) Tag() string {
	return fe.fieldLevel.GetTag()
}

func (fe FieldError) ActualTag() string {
	return fe.fieldLevel.GetTag()
}

func (fe FieldError) Namespace() string {
	return fe.fieldLevel.Parent().Type().Name() + "." + fe.Field()
}

func (fe FieldError) StructNamespace() string {
	return fe.Namespace()
}

func (fe FieldError) Field() string {
	return fe.fieldLevel.FieldName()
}

func (fe FieldError) StructField() string {
	return fe.fieldLevel.StructFieldName()
}

func (fe FieldError) Value() any {
	return fe.val
}

func (fe FieldError) Param() string {
	switch reflect.TypeOf(fe.param).Kind() {
	case reflect.String:
		return fe.param.(string)

	case reflect.Slice:
		if strArr, ok := fe.param.([]string); ok {
			return "[" + strings.Join(strArr, ",") + "]"
		}
	}

	return ""
}

func (fe FieldError) Kind() reflect.Kind {
	return fe.fieldLevel.Field().Kind()
}

func (fe FieldError) Type() reflect.Type {
	return fe.fieldLevel.Field().Type()
}

func (fe FieldError) Error() string {
	return fmt.Sprintf(
		"Key: '%s' Error:Field validation for '%s' failed on the '%s' tag",
		fe.Namespace(),
		fe.Field(),
		fe.Tag(),
	)
}
