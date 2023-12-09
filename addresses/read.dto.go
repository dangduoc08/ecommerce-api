package addresses

import (
	"fmt"

	"github.com/dangduoc08/gooh"
	"github.com/dangduoc08/gooh/common"
	"github.com/dangduoc08/gooh/exception"
	"github.com/go-playground/validator/v10"
)

type READ_Query struct {
	StoreID uint `bind:"store_id" validate:"required"`
}

func (self READ_Query) Transform(query gooh.Query, medata common.ArgumentMetadata) any {
	validate := validator.New()
	queryDTO := query.Bind(self)
	err := validate.Struct(queryDTO)
	errMsgs := []string{}
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", err.Field(), err.Tag()))
		}

		panic(exception.UnprocessableEntityException(errMsgs))
	}

	return queryDTO
}
