package stores

import (
	"github.com/dangduoc08/ecommerce-api/stores/controllers"
	"github.com/dangduoc08/ecommerce-api/stores/providers"
	"github.com/dangduoc08/gooh/core"
)

var Module = core.ModuleBuilder().
	Controllers(
		controllers.StoreREST{},
	).
	Providers(
		providers.StoreDB{},
	).
	Build()
