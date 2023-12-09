package stores

import (
	"fmt"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_ID_Body_Data_DTO struct {
	Name        string `bind:"name" validate:"gte=0,lte=130"`
	Description string `bind:"description"`
	Phone       string `bind:"phone" validate:"phone"`
	Email       string `bind:"email" validate:"omitempty,email"`
}

type UPDATE_BY_ID_Body_DTO struct {
	Data UPDATE_BY_ID_Body_Data_DTO `bind:"data"`
}

func (self UPDATE_BY_ID_Body_DTO) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	errMsgs := []string{}

	validate := validator.New()
	validate.RegisterValidation("phone", utils.ValidatePhone("VN", func(err error) {
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
