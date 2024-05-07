package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type MODIFY_reset_password_Body_Data_DTO struct {
	UserIdentity string `bind:"user_identity" validate:"required"`
}

type MODIFY_reset_password_Body_DTO struct {
	Data MODIFY_reset_password_Body_Data_DTO `bind:"data"`
}

func (instance MODIFY_reset_password_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(instance)
	bodyDTO := bindedStruct.(MODIFY_reset_password_Body_DTO)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	fieldErrs := validate.Struct(bodyDTO)
	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.NestedTag(),
				"error": strings.TrimSpace(fmt.Sprintf("must be %v %v", fieldErr.Tag(), fieldErr.Param())),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}

type MODIFY_reset_password_Header_DTO struct {
	Origin string `bind:"Origin" validate:"required"`
}

func (instance MODIFY_reset_password_Header_DTO) Transform(header gogo.Header, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := header.Bind(instance)
	headerDTO := bindedStruct.(MODIFY_reset_password_Header_DTO)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	fieldErrs := validate.Struct(headerDTO)
	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.NestedTag(),
				"error": strings.TrimSpace(fmt.Sprintf("must be %v %v", fieldErr.Tag(), fieldErr.Param())),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return headerDTO
}
