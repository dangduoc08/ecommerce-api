package conf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/ctx"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Port   int    `bind:"PORT" validate:"required"`
	AppENV string `bind:"APP_ENV" validate:"required,AppENVEnum"`

	Username  string `bind:"USERNAME" validate:"required"`
	Password  string `bind:"PASSWORD" validate:"required"`
	Email     string `bind:"EMAIL" validate:"required,email"`
	FirstName string `bind:"FIRST_NAME" validate:"required"`
	LastName  string `bind:"LAST_NAME" validate:"required"`

	JWTAccessAPIKey      string `bind:"JWT_ACCESS_API_KEY" validate:"required"`
	JWTAccessAPIExpIn    int    `bind:"JWT_ACCESS_API_EXP_IN" validate:"required"`
	JWTRefreshTokenKey   string `bind:"JWT_REFRESH_TOKEN_KEY" validate:"required"`
	JWTRefreshTokenExpIn int    `bind:"JWT_REFRESH_TOKEN_EXP_IN" validate:"required"`

	PostgresHost     string `bind:"POSTGRES_HOST" validate:"required"`
	PostgresUser     string `bind:"POSTGRES_USER" validate:"required"`
	PostgresPassword string `bind:"POSTGRES_PASSWORD" validate:"required"`
	PostgresDB       string `bind:"POSTGRES_DB" validate:"required"`
	PostgresPort     int    `bind:"POSTGRES_PORT" validate:"required"`
	PostgresSSL      bool   `bind:"POSTGRES_SSL" validate:"boolean"`
}

var Value Config

var Module = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {

			// transform to proper types
			dto, fieldLevels := c.Transform(Config{})

			dtoConfig := dto.(Config)
			Value = dtoConfig
			errMsgs := []string{}

			// validate config values should be added correctly
			v := validator.New()
			v.RegisterValidation("AppENVEnum", utils.ValidateEnum(
				constants.APP_ENVS,
				func(err validator.FieldError) {
					if err != nil {
						errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", err.Field(), err.Tag()))
					}
				},
			))

			err := v.Struct(dtoConfig)
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
			dtoConfigType := reflect.TypeOf(dtoConfig)
			for i := 0; i < dtoConfigType.NumField(); i++ {
				field := dtoConfigType.Field(i)
				fieldValue := reflect.ValueOf(dtoConfig).Field(i)
				envKey := field.Tag.Get("bind")
				if envKey != "" {
					c.Set(envKey, fieldValue.Interface())
				}
			}
		},
	},
})
