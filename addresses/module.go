package addresses

import (
	"github.com/dangduoc08/ecommerce-api/addresses/controllers"
	"github.com/dangduoc08/ecommerce-api/addresses/providers"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Controllers(
		controllers.REST{},
	).
	Providers(
		providers.DBHandler{},
	).
	Build()
