package categories

import (
	"github.com/dangduoc08/ecommerce-api/categories/controllers"
	"github.com/dangduoc08/ecommerce-api/categories/providers"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Providers(
		providers.DBHandler{},
		providers.DBValidation{},
	).
	Controllers(
		controllers.REST{},
	).
	Build()
