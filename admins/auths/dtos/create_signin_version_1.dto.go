package dtos

import (
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_signin_VERSION_1_Body_Data_DTO struct {
	Username string `bind:"username" validate:"required,gte=6"`
	Password string `bind:"password" validate:"required,password"`
}

type CREATE_signin_VERSION_1_Body_DTO struct {
	Data CREATE_signin_VERSION_1_Body_Data_DTO `bind:"data"`
}

func (instance CREATE_signin_VERSION_1_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	bindedStruct, fls := body.Bind(instance)
	bodyDTO := bindedStruct.(CREATE_signin_VERSION_1_Body_DTO)

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
				"reason": []string{
					"passwordError",
				},
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
				"reason": []string{
					"mustBe",
					fieldErr.Tag(),
					fieldErr.Param(),
					"characters",
				},
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
