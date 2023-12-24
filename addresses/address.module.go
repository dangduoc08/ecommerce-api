package addresses

import (
	"github.com/dangduoc08/ecommerce-api/addresses/controllers"
	"github.com/dangduoc08/ecommerce-api/addresses/providers"
	"github.com/dangduoc08/gooh/core"
)

var Module = core.ModuleBuilder().
	Controllers(
		controllers.AddressREST{},
	).
	Providers(
		providers.AddressDB{},
	).
	Build()
