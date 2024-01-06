package confs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dangduoc08/ecommerce-api/confs/models"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/ecommerce-api/validators"
	"github.com/dangduoc08/gooh/ctx"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/go-playground/validator/v10"
)

var Value models.Config

var Module = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {

			// transform to proper types
			dto, fieldLevels := c.Transform(models.Config{})

			configDTO := dto.(models.Config)
			Value = configDTO
			errMsgs := []string{}

			// validate config values should be added correctly
			v := validator.New()
			v.RegisterValidation("AppENV", validators.ValidateEnum(
				constants.APP_ENVS,
				func(err validator.FieldError) {
					if err != nil {
						errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", err.Field(), err.Tag()))
					}
				},
			))

			err := v.Struct(configDTO)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					fieldLevel := utils.ArrFind[ctx.FieldLevel](fieldLevels, func(el ctx.FieldLevel, i int) bool {
						return el.Field() == err.Field()
					})

					errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", fieldLevel.Tag(), err.Tag()))
				}
			}

			if len(errMsgs) > 0 {
				panic(strings.Join(errMsgs, "\n       "))
			}

			// re-assign to config struct
			dtoConfigType := reflect.TypeOf(configDTO)
			for i := 0; i < dtoConfigType.NumField(); i++ {
				field := dtoConfigType.Field(i)
				fieldValue := reflect.ValueOf(configDTO).Field(i)
				envKey := field.Tag.Get("bind")
				if envKey != "" {
					c.Set(envKey, fieldValue.Interface())
				}
			}
		},
	},
})
