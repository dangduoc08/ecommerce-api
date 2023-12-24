package permissions

import (
	"github.com/dangduoc08/ecommerce-api/permissions/controllers"
	"github.com/dangduoc08/gooh/core"
)

var Module = core.ModuleBuilder().
	Controllers(
		controllers.PermissionREST{},
	).
	Build()
