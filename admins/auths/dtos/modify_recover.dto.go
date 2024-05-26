package dtos

import (
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type MODIFY_recover_Body_Data_DTO struct {
	Password string `bind:"password" validate:"required,password"`
}

type MODIFY_recover_Body_DTO struct {
	Data MODIFY_recover_Body_Data_DTO `bind:"data"`
}

func (instance MODIFY_recover_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(instance)
	bodyDTO := bindedStruct.(MODIFY_recover_Body_DTO)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	validate.RegisterValidation("password", validators.ValidatePassword(func(fieldErr validator.FieldError) {
		if fieldErr != nil {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field":     fl.Tag(),
				"namespace": fl.NestedTag(),
				"reason": utils.Reason(
					"passwordError",
				),
			})
		}
	}))

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
					fieldErr.Param(),
					"characters",
				),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
