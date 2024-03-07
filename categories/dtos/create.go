package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_Body_Data struct {
	Name              string `bind:"name" validate:"required,gte=1"`
	Description       string `bind:"description"`
	MetaTitle         string `bind:"meta_title" validate:"required,gte=1"`
	MetaDescription   string `bind:"meta_description" validate:"lte=160"`
	Slug              string `bind:"slug" validate:"required,gte=1,slug"`
	Status            string `bind:"status" validate:"required,categoryStatus"`
	ParentCategoryIDs []uint `bind:"parent_category_ids"`
}

type CREATE_Body struct {
	Data CREATE_Body_Data `bind:"data"`
}

func (instance CREATE_Body) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}
	validate := validator.New()

	bindedBody, fls := body.Bind(instance)
	bodyDTO := bindedBody.(CREATE_Body)

	bodyDTO.Data.Name = strings.TrimSpace(bodyDTO.Data.Name)
	bodyDTO.Data.MetaTitle = strings.TrimSpace(bodyDTO.Data.MetaTitle)
	bodyDTO.Data.Slug = strings.TrimSpace(bodyDTO.Data.Slug)
	bodyDTO.Data.ParentCategoryIDs = utils.ArrToUnique(bodyDTO.Data.ParentCategoryIDs)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("categoryStatus", validators.ValidateEnum(constants.CATEGORY_STATUSES, func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v is invalid status", fieldErr.Value()),
			})
		}
	}))

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
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("must be %v %v", fieldErr.Tag(), fieldErr.Param()),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
