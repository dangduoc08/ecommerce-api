package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/assets/providers"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_dirs_Body_Data struct {
	Name string `bind:"name" validate:"required"`
}

type CREATE_dirs_Body struct {
	providers.HandleAsset
	Data CREATE_dirs_Body_Data `bind:"data"`
}

func (self CREATE_dirs_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(self)
	bodyDTO := dto.(CREATE_dirs_Body)
	bodyDTO.Data.Name = strings.TrimSpace(bodyDTO.Data.Name)

	fieldMap := make(map[string]gooh.FieldLevel)
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

	bodyDTO.Data.Name = self.CleanDir(bodyDTO.Data.Name)

	return bodyDTO
}
