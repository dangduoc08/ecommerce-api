package manufacturers

import (
	"github.com/dangduoc08/gogo/core"
)

var ManufacturersModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			ManufacturerController{},
		).
		Providers(
			ManufacturerProvider{},
		).
		Build()

	module.
		Prefix("manufacturers")

	return module
}()
