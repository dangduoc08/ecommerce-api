package groups

import (
	"github.com/dangduoc08/ecommerce-api/groups/controllers"
	"github.com/dangduoc08/ecommerce-api/groups/providers"
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
