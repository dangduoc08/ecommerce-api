package config

import (
	"reflect"

	"github.com/dangduoc08/gooh/modules/config"
)

type Config struct {
	Port         int  `bind:"PORT"`
	PostgresPort int  `bind:"POSTGRES_PORT"`
	PostgresSSL  bool `bind:"POSTGRES_SSL"`
}

var Module = config.Register(&config.ConfigModuleOptions{
	IsGlobal:          true,
	IsExpandVariables: true,
	Hooks: []config.ConfigHookFn{
		func(c config.ConfigService) {
			dtoConfig := c.Transform(Config{}).(Config)
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
