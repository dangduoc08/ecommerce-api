package manufacturers

import (
	"github.com/dangduoc08/ecommerce-api/manufacturers/controllers"
	"github.com/dangduoc08/ecommerce-api/manufacturers/providers"
	"github.com/dangduoc08/gooh/core"
)

var Module = core.ModuleBuilder().
	Controllers(
		controllers.REST{},
	).
	Providers(
		providers.DBHandler{},
	).
	Build()
