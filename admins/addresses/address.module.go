package addresses

import (
	"github.com/dangduoc08/gogo/core"
)

var AddressModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			AddressController{},
		).
		Providers(
			AddressProvider{},
		).
		Build()

	module.
		Prefix("addresses")

	return module
}()
