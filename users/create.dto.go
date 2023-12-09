package users

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/field"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_Body_Data_DTO struct {
	Username  string `bind:"username" validate:"required"`
	Password  string `bind:"password" validate:"required"`
	Email     string `bind:"email" validate:"required,email"`
	FirstName string `bind:"first_name" validate:"required"`
	LastName  string `bind:"last_name" validate:"required"`
}

type CREATE_Body_DTO struct {
	Data         CREATE_Body_Data_DTO `bind:"data"`
	UserProvider UserProvider
}

func (dto CREATE_Body_DTO) Transform(
	body gooh.Body,
	medata common.ArgumentMetadata,
) any {
	validate := validator.New()
	createBody := body.Bind(dto)

	err := validate.Struct(createBody)
	errMsgs := []string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			f := field.UserMapField.GetField(err.Field())
			errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", f, err.Tag()))
		}

		panic(exception.UnprocessableEntityException(errMsgs))
	}

	data := []map[string]string{
		{
			"email": body.Get("data.email").(string),
		},
		{
			"username": body.Get("data.username").(string),
		},
	}

	err = dto.UserProvider.CheckDuplicateUser(data)
	if err != nil {
		errMsgs := []string{
			"Field: user.email or user.username",
			"Error: duplicate data",
		}
		panic(exception.ConflictException(errMsgs))
	}

	return createBody
}
