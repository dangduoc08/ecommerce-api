package products

import (
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/gogo/core"
)

var ProductModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			ProductController{},
		).
		Providers(
			ProductProvider{},
			categories.CategoryProvider{},
		).
		Build()

	module.
		Prefix("products")

	return module
}()
