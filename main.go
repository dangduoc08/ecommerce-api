package main

import (
	"github.com/dangduoc08/ecommerce-api/addresses"
	"github.com/dangduoc08/ecommerce-api/conf"
	"github.com/dangduoc08/ecommerce-api/db"
	"github.com/dangduoc08/ecommerce-api/files"
	"github.com/dangduoc08/ecommerce-api/globals"
	"github.com/dangduoc08/ecommerce-api/groups"
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/ecommerce-api/seeds"
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
		BindGlobalInterceptors(globals.LoggingInterceptor{}, globals.ResponseInterceptor{})

	app.Create(
		core.ModuleBuilder().
			Imports(
				db.DBModule,
				conf.ConfigModule,
				users.UserModule,
				stores.StoreModule,
				seeds.SeedModule,
				locations.LocationModule,
				addresses.AddressModule,
				files.FileModule,
				groups.GroupModule,
			).
			Build(),
	)

	configService := app.Get(config.ConfigService{}).(config.ConfigService)

	app.Logger.Fatal("AppError", "error", app.Listen(configService.Get("PORT").(int)))
}
