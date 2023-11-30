package config

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dangduoc08/gooh/modules/config"
	"github.com/go-playground/validator/v10"
)

type Config struct {
	Port int `bind:"PORT" validate:"required"`

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

var Module = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {

			// transform to proper types
			dtoConfig := c.Transform(Config{}).(Config)

			// validate config values should be added correctly
			validate := validator.New()
			err := validate.Struct(dtoConfig)
			errMsgs := []string{}
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					errMsgs = append(errMsgs, fmt.Sprintf("'%s' %s", err.Field(), err.Tag()))
				}

				panic(strings.Join(errMsgs, ", "))
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
