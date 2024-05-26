package dtos

import (
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/dangduoc08/gogo/modules/config"
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
				"field":     fl.Tag(),
				"namespace": fl.NestedTag(),
				"reason": utils.Reason(
					"mustBe",
					fieldErr.Tag(),
				),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}

type MODIFY_reset_password_Header_DTO struct {
	config.ConfigService
	Origin string `bind:"Origin" validate:"required,cors"`
}

func (instance MODIFY_reset_password_Header_DTO) Transform(header gogo.Header, medata common.ArgumentMetadata) any {
	domainWhitelist := instance.Get("DOMAIN_WHITELIST").([]string)

	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := header.Bind(instance)
	headerDTO := bindedStruct.(MODIFY_reset_password_Header_DTO)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("cors", validators.ValidateEnum(domainWhitelist, func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field":     fl.Tag(),
				"namespace": fl.NestedTag(),
				"reason": utils.Reason(
					"domain",
					fieldErr.Value().(string),
					"nin",
					"availableList",
				),
			})
		}
	}))

	fieldErrs := validate.Struct(headerDTO)
	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field":     fl.Tag(),
				"namespace": fl.NestedTag(),
				"reason": utils.Reason(
					"mustBe",
					fieldErr.Tag(),
				),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return headerDTO
}
