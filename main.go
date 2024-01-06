package main

import (
	"github.com/dangduoc08/ecommerce-api/addresses"
	"github.com/dangduoc08/ecommerce-api/assets"
	"github.com/dangduoc08/ecommerce-api/auths"
	"github.com/dangduoc08/ecommerce-api/categories"
	"github.com/dangduoc08/ecommerce-api/conf"
	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/ecommerce-api/groups"
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/ecommerce-api/permissions"
	"github.com/dangduoc08/ecommerce-api/seeds"
	"github.com/dangduoc08/ecommerce-api/shared"
	"github.com/dangduoc08/ecommerce-api/stores"
	"github.com/dangduoc08/ecommerce-api/users"
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
		BindGlobalInterceptors(shared.LoggingInterceptor{}, shared.ResponseInterceptor{})

	app.Create(
		core.ModuleBuilder().
			Imports(
				db.Module,
				conf.Module,
				auths.Module,
				users.Module,
				stores.Module,
				seeds.Module,
				locations.Module,
				addresses.Module,
				assets.Module,
				groups.Module,
				permissions.Module,
				categories.Module,
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "error", app.Listen(configService.Get("PORT").(int)))
}
