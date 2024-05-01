package stores

import (
	"github.com/dangduoc08/ecommerce-api/admins/addresses"
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	"github.com/dangduoc08/gogo/core"
)

var StoreModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			StoreController{},
		).
		Providers(
			stores.StoreProvider{},
			addresses.AddressProvider{},
			categories.CategoryProvider{},
		).
		Build()

	module.
		Prefix("stores")

	return module
}()
