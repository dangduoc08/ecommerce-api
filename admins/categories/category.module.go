package categories

import (
	"github.com/dangduoc08/gogo/core"
)

var CategoryModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			CategoryController{},
		).
		Providers(
			CategoryProvider{},
		).
		Build()

	module.
		Prefix("categories")

	return module
}()
