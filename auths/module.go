package auths

import (
	"github.com/dangduoc08/ecommerce-api/auths/controllers"
	"github.com/dangduoc08/ecommerce-api/auths/providers"
	"github.com/dangduoc08/ecommerce-api/groups"
	"github.com/dangduoc08/gooh/core"
)

var Module = core.ModuleBuilder().
	Imports(
		groups.Module,
	).
	Controllers(
		controllers.REST{},
	).
	Providers(
		providers.DBHandler{},
	).
	Build()
