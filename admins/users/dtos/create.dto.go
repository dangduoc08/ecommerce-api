package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_Body_Data_DTO struct {
	Username  string `bind:"username" validate:"required,gte=6,username"`
	Password  string `bind:"password" validate:"required,password"`
	Email     string `bind:"email" validate:"required,email"`
	FirstName string `bind:"first_name" validate:"required"`
	LastName  string `bind:"last_name" validate:"required"`
	GroupIDs  []uint `bind:"group_ids"`
}

type CREATE_Body_DTO struct {
	Data CREATE_Body_Data_DTO `bind:"data"`
}

func (instance CREATE_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(instance)
	bodyDTO := bindedStruct.(CREATE_Body_DTO)

	bodyDTO.Data.Username = strings.TrimSpace(bodyDTO.Data.Username)
	bodyDTO.Data.FirstName = strings.TrimSpace(bodyDTO.Data.FirstName)
	bodyDTO.Data.FirstName = strings.TrimSpace(bodyDTO.Data.FirstName)
	bodyDTO.Data.Email = strings.TrimSpace(bodyDTO.Data.Email)
	bodyDTO.Data.GroupIDs = utils.ArrToUnique(bodyDTO.Data.GroupIDs)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("username", validators.ValidateUsername(func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": "must include valid characters",
			})
		}
	}))

	validate.RegisterValidation("password", validators.ValidatePassword(func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": "must be at least 8 characters including 1 upper case, 1 digit and 1 special character",
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
