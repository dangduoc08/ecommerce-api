package stores

import (
	"github.com/dangduoc08/ecommerce-api/addresses"
	"github.com/dangduoc08/ecommerce-api/categories"
	"github.com/dangduoc08/ecommerce-api/stores/controllers"
	"github.com/dangduoc08/ecommerce-api/stores/providers"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Imports(
		categories.Module,
		addresses.Module,
	).
	Controllers(
		controllers.REST{},
	).
	Providers(
		providers.DBHandler{},
	).
	Build()
