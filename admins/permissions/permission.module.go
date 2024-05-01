package permissions

import (
	"github.com/dangduoc08/gogo/core"
)

var PermissionModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			PermissionController{},
		).
		Build()

	module.
		Prefix("permissions")

	return module
}()
