package dtos

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_sessions_Body_Data struct {
	Username string `bind:"username" validate:"required,gte=6"`
	Password string `bind:"password" validate:"required,password"`
}

type CREATE_sessions_Body struct {
	Data CREATE_sessions_Body_Data `bind:"data"`
}

func (self CREATE_sessions_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(self)
	bodyDTO := bindedStruct.(CREATE_sessions_Body)

	fieldMap := make(map[string]gooh.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("password", validators.ValidatePassword(func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprint("must be at least 8 characters including 1 upper case, 1 digit and 1 special character"),
			})
		}
	}))

	fieldErrs := validate.Struct(bodyDTO)
	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("must be %v %v", fieldErr.Tag(), fieldErr.Param()),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
