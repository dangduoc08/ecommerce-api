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
	Name              string `bind:"name" validate:"required,gte=5"`
	Description       string `bind:"description"`
	MetaTitle         string `bind:"meta_title" validate:"required,gte=5"`
	MetaDescription   string `bind:"meta_description" validate:"lte=160"`
	Slug              string `bind:"slug" validate:"required,gte=1,slug"`
	Status            string `bind:"status" validate:"required,categoryStatus"`
	ParentCategoryIDs []uint `bind:"parent_category_ids"`
}

type UPDATE_BY_id_Body_DTO struct {
	Data UPDATE_BY_id_Body_Data_DTO `bind:"data"`
}

func (instance UPDATE_BY_id_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}
	validate := validator.New()

	bindedBody, fls := body.Bind(instance)
	bodyDTO := bindedBody.(UPDATE_BY_id_Body_DTO)

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
