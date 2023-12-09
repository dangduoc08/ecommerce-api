package conf

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dangduoc08/ecommerce-api/utils"
	"github.com/dangduoc08/gooh/modules/config"
	"github.com/go-playground/validator/v10"
)

const (
	APP_ENV_DEV  = "development"
	APP_ENV_TEST = "test"
	APP_ENV_PROD = "production"
)

type Config struct {
	Port   int    `bind:"PORT" validate:"required"`
	AppENV string `bind:"APP_ENV" validate:"AppENVEnum"`

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

var ConfigModule = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {

			// transform to proper types
			dtoConfig := c.Transform(Config{}).(Config)
			Value = dtoConfig
			errMsgs := []string{}

			// validate config values should be added correctly
			v := validator.New()
			v.RegisterValidation("AppENVEnum", utils.ValidateEnum(
				[]string{APP_ENV_DEV, APP_ENV_TEST, APP_ENV_PROD},
				func(err error) {
					if err != nil {
						errMsgs = append(errMsgs, err.Error())
					}
				},
			))

			err := v.Struct(dtoConfig)
			if err != nil {
				for _, err := range err.(validator.ValidationErrors) {
					errMsgs = append(errMsgs, fmt.Sprintf("Field: %s, Error: must be %s", err.Field(), err.Tag()))
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
