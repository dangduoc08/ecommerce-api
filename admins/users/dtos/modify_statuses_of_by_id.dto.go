package dtos

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type MODIFY_statuses_OF_BY_id_Param_DTO struct {
	ID uint `bind:"id" validate:"required"`
}

func (instance MODIFY_statuses_OF_BY_id_Param_DTO) Transform(param gogo.Param, medata common.ArgumentMetadata) any {
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

type MODIFY_statuses_OF_BY_id_Data_DTO struct {
	Status string `bind:"status" validate:"userStatus"`
}

type MODIFY_statuses_OF_BY_id_DTO struct {
	Data MODIFY_statuses_OF_BY_id_Data_DTO `bind:"data"`
}

func (instance MODIFY_statuses_OF_BY_id_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(instance)
	bodyDTO := bindedStruct.(MODIFY_statuses_OF_BY_id_DTO)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("userStatus", validators.ValidateEnum(constants.USER_STATUSES, func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v is invalid status", fieldErr.Value()),
			})
		}
	}))

	validate.Struct(bodyDTO)

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
