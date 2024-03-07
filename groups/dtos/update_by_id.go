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

type UPDATE_BY_id_Param struct {
	ID uint `bind:"id"`
}

func (instance UPDATE_BY_id_Param) Transform(param gogo.Param, medata common.ArgumentMetadata) any {
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

type UPDATE_BY_id_Body_Data struct {
	Name        string   `bind:"name" validate:"required"`
	Permissions []string `bind:"permissions" validate:"permissions"`
}

type UPDATE_BY_id_Body struct {
	common.REST
	Data UPDATE_BY_id_Body_Data `bind:"data"`
}

func (instance UPDATE_BY_id_Body) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(instance)
	bodyDTO := dto.(UPDATE_BY_id_Body)

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
				"error": fmt.Sprintf("must be %v", fieldErr.Tag()),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	bodyDTO.Data.Permissions = utils.ArrToUnique[string](bodyDTO.Data.Permissions)

	return bodyDTO
}
