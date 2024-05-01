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

type READ_Query_DTO struct {
	Sort   string `bind:"sort"`
	Order  string `bind:"order" validate:"order"`
	Limit  int    `bind:"limit" validate:"gte=0,lte=100"`
	Offset int    `bind:"offset" validate:"gte=0"`
}

func (instance READ_Query_DTO) Transform(query gogo.Query, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := query.Bind(instance)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	queryDTO := bindedStruct.(READ_Query_DTO)

	validate.RegisterValidation("order", validators.ValidateEnum(constants.ORDERS, func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v not in %v", fieldErr.Value(), fieldErr.Param()),
			})
		}
	}))

	fieldErrs := validate.Struct(queryDTO)

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

	return queryDTO
}
