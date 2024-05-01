package stores

import (
	"github.com/dangduoc08/gogo/core"
)

var StoreModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			StoreController{},
		).
		Providers(
			StoreProvider{},
		).
		Build()

	module.
		Prefix("stores")

	return module
}()
