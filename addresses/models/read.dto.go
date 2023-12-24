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

type READ_Query struct {
	StoreID uint   `bind:"store_id" validate:"required"`
	Sort    string `bind:"sort"`
	Order   string `bind:"order" validate:"order"`
	Limit   int    `bind:"limit" validate:"gte=0,lte=100"`
	Offset  int    `bind:"offset" validate:"gte=0"`
}

func (self READ_Query) Transform(query gooh.Query, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := query.Bind(self)
	queryDTO := bindedStruct.(READ_Query)

	fieldMap := make(map[string]gooh.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("order", utils.ValidateEnum(constants.ORDERS, func(fieldErr validator.FieldError) {
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