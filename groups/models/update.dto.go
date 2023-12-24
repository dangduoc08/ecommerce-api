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
	Name        string   `bind:"name" validate:"required"`
	Permissions []string `bind:"permissions" validate:"permissions"`
}

type UPDATE_BY_id_Body struct {
	common.REST
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

	availablePermissions := []string{}
	for _, restConfiguration := range self.GetConfigurations() {
		availablePermissions = append(
			availablePermissions,
			restConfiguration.Method+restConfiguration.Route,
		)
	}

	validate.RegisterValidation("permissions", utils.ValidateEnums(availablePermissions, func(fieldErr validator.FieldError) {
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
