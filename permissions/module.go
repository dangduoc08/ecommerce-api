package permissions

import (
	"github.com/dangduoc08/ecommerce-api/permissions/controllers"
	"github.com/dangduoc08/gogo/core"
)

var Module = core.ModuleBuilder().
	Controllers(
		controllers.REST{},
	).
	Build()
