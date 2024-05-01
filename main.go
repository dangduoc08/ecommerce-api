package main

import (
	"github.com/dangduoc08/ecommerce-api/admins"
	"github.com/dangduoc08/ecommerce-api/confs"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/ecommerce-api/seeds"
	sharedLayers "github.com/dangduoc08/ecommerce-api/shared_layers"
	"github.com/dangduoc08/ecommerce-api/statics"
	"github.com/dangduoc08/ecommerce-api/storefronts"
	"github.com/dangduoc08/gogo/core"
	"github.com/dangduoc08/gogo/log"
	"github.com/dangduoc08/gogo/middlewares"
	"github.com/dangduoc08/gogo/modules/config"
	"github.com/dangduoc08/gogo/versioning"
)

func main() {
	app := core.New()

	logger := log.NewLog(&log.LogOptions{
		Level:     log.DebugLevel,
		LogFormat: log.PrettyFormat,
	})

	app.
		EnableVersioning(
			versioning.Versioning{
				Type:           versioning.HEADER,
				Key:            "v",
				DefaultVersion: versioning.NEUTRAL_VERSION,
			},
		).
		UseLogger(logger).
		Use(
			middlewares.CORS(),
			middlewares.RequestLogger(logger),
		).
		BindGlobalInterceptors(
			sharedLayers.LoggingInterceptor{},
			sharedLayers.ResponseInterceptor{},
		)

	app.Create(
		core.ModuleBuilder().
			Imports(
				dbs.DBModule,
				confs.ConfModule,
				seeds.SeedModule,
				admins.AdminModule,
				storefronts.StorefrontModule,
				statics.StaticModule,
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "error", app.Listen(configService.Get("PORT").(int)))
}
