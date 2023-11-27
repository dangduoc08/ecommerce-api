package main

import (
	"strconv"

	"github.com/dangduoc08/gooh/core"
	"github.com/dangduoc08/gooh/log"
	"github.com/dangduoc08/gooh/middlewares"
	"github.com/dangduoc08/gooh/modules/config"
)

func main() {
	app := core.New()
	logger := log.NewLog(&log.LogOptions{
		Level:     log.DebugLevel,
		LogFormat: log.PrettyFormat,
	})

	app.
		UseLogger(logger).
		Use(middlewares.CORS(), middlewares.RequestLogger(logger))

	app.Create(
		core.ModuleBuilder().
			Imports(
				config.Register(&config.ConfigModuleOptions{
					IsGlobal:          true,
					IsExpandVariables: true,
					Hooks: []config.ConfigHookFn{
						func(c config.ConfigService) {
							port := c.Get("PORT")
							if s, ok := port.(string); ok {
								port, err := strconv.Atoi(s)
								if err != nil {
									panic(err)
								}
								c.Set("PORT", port)
							}
						},
					},
				}),
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "errMsg", app.Listen(configService.Get("PORT").(int)))
}