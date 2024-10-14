package main

import (
	"github.com/dangduoc08/ecommerce-api/admins"
	"github.com/dangduoc08/ecommerce-api/confs"
	"github.com/dangduoc08/ecommerce-api/dbs"
	"github.com/dangduoc08/ecommerce-api/mails"
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
				DefaultVersion: versioning.NEUTRAL_VERSION,
				Key:            "X-Api-Version",
			},
		).
		UseLogger(logger).
		Use(
			middlewares.RequestLogger(logger),
			middlewares.CORS(&middlewares.CORSOptions{
				AllowOrigin:        confs.ENV.DomainWhitelist,
				IsAllowCredentials: true,
				AllowHeaders: []string{
					"Origin",
					"X-Requested-With",
					"Content-Type",
					"Accept",
					"v",
					"Cookie",
				},
			}),
		).
		BindGlobalInterceptors(
			sharedLayers.LoggingInterceptor{},
			sharedLayers.ResponseInterceptor{},
		).
		BindGlobalExceptionFilters(
			sharedLayers.AllExceptionsFilter{},
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
				mails.MailModule,
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "error", app.Listen(configService.Get("PORT").(int)))
}
