package products

import (
	"github.com/dangduoc08/ecommerce-api/categories"
	"github.com/dangduoc08/ecommerce-api/products/controllers"
	"github.com/dangduoc08/ecommerce-api/products/providers"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Imports(
		categories.Module,
	).
	Controllers(
		controllers.REST{},
	).
	Providers(
		providers.DBHandler{},
	).
	Build()
