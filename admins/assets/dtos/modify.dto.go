package dtos

import (
	"fmt"
	"strings"

	"github.com/dangduoc08/ecommerce-api/admins/assets/commons"
	"github.com/dangduoc08/gogo"
	"github.com/dangduoc08/gogo/common"
	"github.com/dangduoc08/gogo/exception"
	"github.com/go-playground/validator/v10"
)

type MODIFY_Body_Data_DTO struct {
	OldDir string `bind:"old_dir" validate:"required"`
	NewDir string `bind:"new_dir" validate:"required"`
}

type MODIFY_Body_DTO struct {
	CommonProvider commons.CommonProvider
	Data           MODIFY_Body_Data_DTO `bind:"data"`
}

func (instance MODIFY_Body_DTO) Transform(body gogo.Body, medata common.ArgumentMetadata) any {
	errMsgs := []map[string]any{}

	validate := validator.New()
	dto, fls := body.Bind(instance)
	bodyDTO := dto.(MODIFY_Body_DTO)
	bodyDTO.Data.NewDir = strings.TrimSpace(bodyDTO.Data.NewDir)

	fieldMap := make(map[string]gogo.FieldLevel)
	for _, fl := range fls {
		fieldMap[fl.Field()] = fl
	}

	fieldErrs := validate.Struct(bodyDTO)

	if fieldErrs != nil {
		for _, fieldErr := range fieldErrs.(validator.ValidationErrors) {
			fl := fieldMap[fieldErr.Field()]
			errMsgs = append(errMsgs, map[string]any{
				"field": fl.Tag(),
				"error": fmt.Sprintf("must be %v", fieldErr.Tag()),
			})
		}
	}

	if len(errMsgs) > 0 {
		panic(exception.UnprocessableEntityException(errMsgs))
	}

	bodyDTO.Data.OldDir = instance.CommonProvider.CleanDir(bodyDTO.Data.OldDir)
	bodyDTO.Data.NewDir = instance.CommonProvider.CleanDir(bodyDTO.Data.NewDir)

	return bodyDTO
}