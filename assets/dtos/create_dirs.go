package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_dirs_Body_Data struct {
	Dir string `bind:"dir" validate:"required"`
}

type CREATE_dirs_Body struct {
	providers.HandleAsset
	Data CREATE_dirs_Body_Data `bind:"data"`
}

func (instance CREATE_dirs_Body) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(instance)
	bodyDTO := dto.(CREATE_dirs_Body)
	bodyDTO.Data.Dir = strings.TrimSpace(bodyDTO.Data.Dir)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	fieldErrs := validate.Struct(bodyDTO)

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

	bodyDTO.Data.Dir = instance.CleanDir(bodyDTO.Data.Dir)

	return bodyDTO
}
