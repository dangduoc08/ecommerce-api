package locations

import (
	"github.com/dangduoc08/ecommerce-api/locations/controllers"
	"github.com/dangduoc08/ecommerce-api/locations/providers"
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
