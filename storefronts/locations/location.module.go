package locations

import (
	"github.com/dangduoc08/ecommerce-api/admins/locations"
	"github.com/dangduoc08/gogo/core"
)

var LocationModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			LocationController{},
		).
		Providers(
			locations.LocationProvider{},
		).
		Build()

	module.
		Prefix("locations")

	return module
}()
