package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param_DTO struct {
	ID uint `bind:"id" validate:"required"`
}

func (instance UPDATE_BY_id_Param_DTO) Transform(param gogo.Param, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	paramDTO, fls := param.Bind(instance)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	fieldErrs := validate.Struct(paramDTO)

	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("must be %v", fieldErr.Tag()),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return paramDTO
}

type UPDATE_BY_id_Body_Data_DTO struct {
	Name        string `bind:"name" validate:"required,lte=130"`
	Description string `bind:"description"`
	Phone       string `bind:"phone" validate:"phone"`
	Email       string `bind:"email" validate:"omitempty,email"`
}

type UPDATE_BY_id_Body_DTO struct {
	Data UPDATE_BY_id_Body_Data_DTO `bind:"data"`
}

func (instance UPDATE_BY_id_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(instance)
	bodyDTO := dto.(UPDATE_BY_id_Body_DTO)

	bodyDTO.Data.Name = strings.TrimSpace(bodyDTO.Data.Name)
	bodyDTO.Data.Description = strings.TrimSpace(bodyDTO.Data.Description)
	bodyDTO.Data.Phone = strings.TrimSpace(bodyDTO.Data.Phone)
	bodyDTO.Data.Email = strings.TrimSpace(bodyDTO.Data.Email)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("phone", validators.ValidatePhone("VN", func(fieldErr validator.FieldError) {
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
