package groups

import (
	"github.com/dangduoc08/gogo/core"
)

var GroupModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			GroupController{},
		).
		Providers(
			GroupProvider{},
		).
		Build()

	module.
		Prefix("groups")

	return module
}()
