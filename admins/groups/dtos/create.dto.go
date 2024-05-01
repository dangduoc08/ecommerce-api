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
	Name        string   `bind:"name" validate:"required"`
	Permissions []string `bind:"permissions" validate:"permissions"`
}

type CREATE_Body_DTO struct {
	common.REST
	Data CREATE_Body_Data_DTO `bind:"data"`
}

func (instance CREATE_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedBody, fls := body.Bind(instance)
	bodyDTO := bindedBody.(CREATE_Body_DTO)

	bodyDTO.Data.Name = strings.TrimSpace(bodyDTO.Data.Name)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	availablePermissions := []string{}
	for _, restConfiguration := range instance.GetConfigurations() {
		availablePermissions = append(
			availablePermissions,
			restConfiguration.Method+restConfiguration.Route,
		)
	}

	validate.RegisterValidation("permissions", validators.ValidateEnums(availablePermissions, func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v is invalid permission", fieldErr.Value()),
			})
		}
	}))

	fieldErrs := validate.Struct(bodyDTO)

	if len(bodyDTO.Data.Permissions) == 1 && bodyDTO.Data.Permissions[0] == "*" {
		errMsgs = []map[string]any{}
	}

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

	bodyDTO.Data.Permissions = utils.ArrToUnique[string](bodyDTO.Data.Permissions)

	return bodyDTO
}
