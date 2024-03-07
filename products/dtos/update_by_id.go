package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/products/models"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_id_Param struct {
	ID uint `bind:"id"`
}

func (instance UPDATE_BY_id_Param) Transform(param gogo.Param, medata common.ArgumentMetadata) any {
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

type UPDATE_BY_id_Body_Data struct {
	Name            string                `bind:"name" validate:"required,gte=1"`
	Description     string                `bind:"description"`
	MetaTitle       string                `bind:"meta_title" validate:"required,gte=1"`
	MetaDescription string                `bind:"meta_description" validate:"lte=160"`
	Slug            string                `bind:"slug" validate:"required,gte=1,slug"`
	Quantity        int                   `bind:"quantity" validate:"gte=0"`
	SKU             string                `bind:"sku"`
	Height          float64               `bind:"height" validate:"gte=0"`
	Width           float64               `bind:"width" validate:"gte=0"`
	Length          float64               `bind:"length" validate:"gte=0"`
	Weight          float64               `bind:"weight" validate:"gte=0"`
	CategoryIDs     []uint                `bind:"category_ids"`
	VariantIDs      []uint                `bind:"variant_ids"`
	ManufacturerID  uint                  `bind:"manufacturer_id"`
	Status          string                `bind:"status"`
	Images          []models.ProductImage `bind:"images" validate:"dive"`
}

type UPDATE_BY_id_Body struct {
	Data UPDATE_BY_id_Body_Data `bind:"data"`
}

func (instance UPDATE_BY_id_Body) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}
	validate := validator.New()

	bindedBody, fls := body.Bind(instance)
	bodyDTO := bindedBody.(UPDATE_BY_id_Body)

	bodyDTO.Data.Name = strings.TrimSpace(bodyDTO.Data.Name)
	bodyDTO.Data.MetaTitle = strings.TrimSpace(bodyDTO.Data.MetaTitle)
	bodyDTO.Data.Slug = strings.TrimSpace(bodyDTO.Data.Slug)
	bodyDTO.Data.SKU = strings.TrimSpace(bodyDTO.Data.SKU)
	bodyDTO.Data.CategoryIDs = utils.ArrToUnique(bodyDTO.Data.CategoryIDs)
	bodyDTO.Data.VariantIDs = utils.ArrToUnique(bodyDTO.Data.VariantIDs)

	fieldMap := make(map[string]gogo.FieldLevel)
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

	validate.RegisterValidation("dir", validators.ValidateDir(func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("%v is invalid dir", fieldErr.Value()),
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
