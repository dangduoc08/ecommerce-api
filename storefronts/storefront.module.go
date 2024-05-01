package storefronts

import (
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/ecommerce-api/storefronts/locations"
	"github.com/dangduoc08/ecommerce-api/storefronts/stores"
	"github.com/dangduoc08/gogo/core"
)

var StorefrontModule = func() *core.Module {
	var module = core.ModuleBuilder().
		Imports(
			locations.LocationModule,
			stores.StoreModule,
		).
		Build()

	module.
		Prefix(constants.STOREFRONT_PATH)

	return module
}
