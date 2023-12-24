package models

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type MODIFY_BY_id_Body_Data struct {
	Status string `bind:"status" validate:"userStatus"`
}

type MODIFY_BY_id_Body struct {
	Data MODIFY_BY_id_Body_Data `bind:"data"`
}

func (self MODIFY_BY_id_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(self)
	bodyDTO := bindedStruct.(MODIFY_BY_id_Body)

	fieldMap := make(map[string]gooh.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("userStatus", utils.ValidateEnum(constants.USER_STATUSES, func(fieldErr validator.FieldError) {
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
