package groups

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type CREATE_Body_Data_DTO struct {
	Name        string   `bind:"name" validate:"required"`
	Permissions []string `bind:"permissions" validate:"mustBeValidPermission"`
}

type CREATE_Body_DTO struct {
	common.REST
	Data CREATE_Body_Data_DTO `bind:"data"`
}

func (self CREATE_Body_DTO) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []string{}
	availablePermissions := []string{}

	for _, restConfiguration := range self.GetConfigurations() {
		availablePermissions = append(
			availablePermissions,
			restConfiguration.Method+restConfiguration.Route,
		)
	}

	validate := validator.New()
	validate.RegisterValidation("mustBeValidPermission", utils.ValidateEnums(availablePermissions, func(err error) {
		if err != nil {
			errMsgs = append(errMsgs, err.Error())
		}
	}))

	bodyDTO := body.Bind(self)

	err := validate.Struct(bodyDTO)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", err.Field(), err.Tag()))
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return bodyDTO
}
