package main

import (
	appConfig "github.com/dangduoc08/ecommerce-api/config"
	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/ecommerce-api/user"
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
		Use(middlewares.CORS(), middlewares.RequestLogger(logger)).
		BindGlobalInterceptors(globals.LoggingInterceptor{}, globals.ResponseInterceptor{})

	app.Create(
		core.ModuleBuilder().
			Imports(
				appConfig.Module,
				user.Module,
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "error", app.Listen(configService.Get("PORT").(int)))
}
