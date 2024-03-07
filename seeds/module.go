package seeds

import (
	"github.com/dangduoc08/ecommerce-api/auths"
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/ecommerce-api/seeds/providers"
	"github.com/dangduoc08/ecommerce-api/stores"
	"github.com/dangduoc08/ecommerce-api/users"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Imports(
		users.Module,
		stores.Module,
		locations.Module,
		auths.Module,
	).
	Providers(
		providers.Seed{},
	).
	Build()
