package addresses

import (
	"fmt"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type UPDATE_BY_ID_Param_DTO struct {
	ID uint `bind:"id" validate:"required"`
}

func (self UPDATE_BY_ID_Param_DTO) Transform(param gooh.Param, medata common.ArgumentMetadata) any {
	validate := validator.New()
	paramDTO := param.Bind(self)
	err := validate.Struct(paramDTO)
	errMsgs := []string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", err.Field(), err.Tag()))
		}

		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return paramDTO
}

type UPDATE_BY_ID_Body_Data_DTO struct {
	StreetName string `bind:"street_name"`
	LocationID uint   `bind:"location_id"`
}

type UPDATE_BY_ID_Body_DTO struct {
	Data UPDATE_BY_ID_Body_Data_DTO `bind:"data"`
}

func (self UPDATE_BY_ID_Body_DTO) Transform(body gooh.Body, medata common.ArgumentMetadata) any {
	return body.Bind(self)
}
