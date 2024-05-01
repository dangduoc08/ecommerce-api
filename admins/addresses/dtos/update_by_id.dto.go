package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param_DTO struct {
	ID uint `bind:"id"`
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
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type UPDATE_BY_id_Body_DTO struct {
	Data UPDATE_BY_id_Body_Data_DTO `bind:"data"`
}

func (instance UPDATE_BY_id_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	bindedBody, _ := body.Bind(instance)

	bodyDTO := bindedBody.(UPDATE_BY_id_Body_DTO)
	bodyDTO.Data.StreetName = strings.TrimSpace(bodyDTO.Data.StreetName)

	return bodyDTO
}
