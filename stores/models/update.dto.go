package models

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param struct {
	ID uint `bind:"id"`
}

func (self UPDATE_BY_id_Param) Transform(param gooh.Param, medata common.ArgumentMetadata) any {
	paramDTO, _ := param.Bind(self)
	return paramDTO
}

type UPDATE_BY_id_Body_Data struct {
	Name        string `bind:"name" validate:"required,lte=130"`
	Description string `bind:"description"`
	Phone       string `bind:"phone" validate:"phone"`
	Email       string `bind:"email" validate:"omitempty,email"`
}

type UPDATE_BY_id_Body struct {
	Data UPDATE_BY_id_Body_Data `bind:"data"`
}

func (self UPDATE_BY_id_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(self)
	bodyDTO := dto.(UPDATE_BY_id_Body)

	fieldMap := make(map[string]gooh.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("phone", utils.ValidatePhone("VN", func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v is invalid phone format", fieldErr.Value()),
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
