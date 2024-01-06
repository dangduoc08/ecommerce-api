package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param struct {
	ID uint `bind:"id" validate:"required"`
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
	Name string `bind:"name" validate:"required,gte=5"`
	Slug string `bind:"slug" validate:"required,gte=1,slug"`
	Logo string `bind:"logo" validate:"omitempty,http_url"`
}

type UPDATE_BY_id_Body struct {
	Data UPDATE_BY_id_Body_Data `bind:"data"`
}

func (self UPDATE_BY_id_Body) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}
	validate := validator.New()

	bindedBody, fls := body.Bind(self)
	bodyDTO := bindedBody.(UPDATE_BY_id_Body)

	bodyDTO.Data.Name = strings.TrimSpace(bodyDTO.Data.Name)
	bodyDTO.Data.Slug = strings.TrimSpace(bodyDTO.Data.Slug)
	bodyDTO.Data.Logo = strings.TrimSpace(bodyDTO.Data.Logo)

	fieldMap := make(map[string]gooh.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("slug", validators.ValidateSlug(func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v is invalid slug", fieldErr.Value()),
			})
		}
	}))

	fieldErrs := validate.Struct(bodyDTO)

	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			if fieldErr.Param() != "" {
				errMsgs = append(errMsgs, map[string]any{
					"field": fl.Tag(),
					"error": fmt.Sprintf("must be %v %v", fieldErr.Tag(), fieldErr.Param()),
				})
			} else {
				errMsgs = append(errMsgs, map[string]any{
					"field": fl.Tag(),
					"error": fmt.Sprintf("must be %v", fieldErr.Tag()),
				})
			}
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
