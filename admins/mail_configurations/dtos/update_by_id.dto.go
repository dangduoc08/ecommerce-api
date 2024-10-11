package dtos

import (
	"strings"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param_DTO struct {
	ID uint `bind:"id" validate:"required"`
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
				"field":     fl.Tag(),
				"namespace": fl.NestedTag(),
				"reason": []string{
					"mustBe",
					fieldErr.Tag(),
					fieldErr.Param(),
				},
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return paramDTO
}

type UPDATE_BY_id_Body_Data_DTO struct {
	Host     string `bind:"host" validate:"omitempty,fqdn"`
	Port     int    `bind:"port"`
	Username string `bind:"username"`
	Password string `bind:"password"`
}

type UPDATE_BY_id_Body_DTO struct {
	Data UPDATE_BY_id_Body_Data_DTO `bind:"data" validate:"required"`
}

func (instance UPDATE_BY_id_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(instance)
	bodyDTO := dto.(UPDATE_BY_id_Body_DTO)

	bodyDTO.Data.Host = strings.TrimSpace(bodyDTO.Data.Host)
	bodyDTO.Data.Username = strings.TrimSpace(bodyDTO.Data.Username)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	fieldErrs := validate.Struct(bodyDTO)

	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field":     fl.Tag(),
				"namespace": fl.NestedTag(),
				"reason": []string{
					"mustBe",
					fieldErr.Tag(),
					fieldErr.Param(),
				},
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
