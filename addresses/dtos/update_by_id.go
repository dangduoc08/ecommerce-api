package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param struct {
	ID uint `bind:"id"`
}

func (self UPDATE_BY_id_Param) Transform(param gooh.Param, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	paramDTO, fls := param.Bind(self)

	fieldMap := make(map[string]gooh.FieldLevel)
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
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type UPDATE_BY_id_Body struct {
	Data UPDATE_BY_id_Body_Data `bind:"data"`
}

func (self UPDATE_BY_id_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	bindedBody, _ := body.Bind(self)

	bodyDTO := bindedBody.(CREATE_Body)
	bodyDTO.Data.StreetName = strings.TrimSpace(bodyDTO.Data.StreetName)

	return bodyDTO
}
