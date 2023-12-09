package users

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/field"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_sessions_Body_Data_DTO struct {
	ID       uint   `bind:"ID"`
	Username string `bind:"username" validate:"required"`
	Password string `bind:"password" validate:"required"`
}

type CREATE_sessions_Body_DTO struct {
	Data CREATE_sessions_Body_Data_DTO `bind:"data"`
}

func (dto CREATE_sessions_Body_DTO) Transform(
	body gooh.Body,
	medata common.ArgumentMetadata,
) any {
	validate := validator.New()
	signinBody := body.Bind(dto)
	err := validate.Struct(signinBody)
	errMsgs := []string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			f := field.UserMapField.GetField(err.Field())
			errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", f, err.Tag()))
		}

		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return signinBody
}
