package addresses

import (
	"github.com/dangduoc08/ecommerce-api/locations"
	"github.com/dangduoc08/gooh/core"
)

var AddressModule = core.ModuleBuilder().
	Imports(
		locations.LocationModule,
	).
	Controllers(
		AddressController{},
	).
	Providers(
		AddressProvider{},
	).
	Exports(
		AddressProvider{},
	).
	Build()
