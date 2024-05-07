package confs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gogo/ctx"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/go-playground/validator/v10"
)

var Value ConfModel

var ConfModule = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {

			// transform to proper types
			dto, fieldLevels := c.Transform(ConfModel{})

			configDTO := dto.(ConfModel)
			Value = configDTO
			errMsgs := []string{}

			// validate config values should be added correctly
			v := validator.New()

			err := v.Struct(configDTO)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					fieldLevel := utils.ArrFind(fieldLevels, func(el ctx.FieldLevel, i int) bool {
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
