package dtos

import (
	"fmt"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_signin_Body_Data_DTO struct {
	Username string `bind:"username" validate:"required"`
	Password string `bind:"password" validate:"required"`
}

type CREATE_signin_Body_DTO struct {
	Data CREATE_signin_Body_Data_DTO `bind:"data"`
}

func (dto CREATE_signin_Body_DTO) Transform(
	body gooh.Body,
	medata common.ArgumentMetadata,
) any {
	validate := validator.New()
	signinBody := body.Bind(dto)

	err := validate.Struct(signinBody)
	errMsgs := []string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: %s", err.Field(), err.Tag()))
		}

		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return signinBody
}
