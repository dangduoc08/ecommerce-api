package seeds

import (
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/ecommerce-api/stores"
	"github.com/dangduoc08/ecommerce-api/users"
	"github.com/dangduoc08/gooh/core"
)

var SeedModule = core.ModuleBuilder().
	Imports(
		users.UserModule,
		stores.StoreModule,
		locations.LocationModule,
	).
	Providers(
		SeedProvider{},
	).
	Build()
