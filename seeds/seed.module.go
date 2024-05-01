package seeds

import (
	"github.com/dangduoc08/ecommerce-api/admins/auths"
	"github.com/dangduoc08/ecommerce-api/storefronts/locations"
	"github.com/dangduoc08/gogo/core"
)

var SeedModule = core.ModuleBuilder().
	Imports(
		locations.LocationModule,
		auths.AuthModule,
	).
	Providers(
		SeedProvider{},
	).
	Build()
