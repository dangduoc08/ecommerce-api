package admins

import (
	"github.com/dangduoc08/ecommerce-api/admins/addresses"
	"github.com/dangduoc08/ecommerce-api/admins/assets"
	"github.com/dangduoc08/ecommerce-api/admins/auths"
	"github.com/dangduoc08/ecommerce-api/admins/categories"
	"github.com/dangduoc08/ecommerce-api/admins/groups"
	mailConfigurations "github.com/dangduoc08/ecommerce-api/admins/mail_configurations"
	"github.com/dangduoc08/ecommerce-api/admins/manufacturers"
	"github.com/dangduoc08/ecommerce-api/admins/permissions"
	"github.com/dangduoc08/ecommerce-api/admins/products"
	"github.com/dangduoc08/ecommerce-api/admins/stores"
	"github.com/dangduoc08/ecommerce-api/admins/users"
	"github.com/dangduoc08/ecommerce-api/constants"
	"github.com/dangduoc08/gogo/core"
)

var AdminModule = func() *core.Module {
	var module = core.ModuleBuilder().
		Imports(
			users.UserModule,
			auths.AuthModule,
			addresses.AddressModule,
			groups.GroupModule,
			permissions.PermissionModule,
			stores.StoreModule,
			manufacturers.ManufacturersModule,
			categories.CategoryModule,
			products.ProductModule,
			assets.AssetModule,
			mailConfigurations.MailConfigurationModule,
		).
		Build()

	module.
		Prefix(constants.ADMIN_PATH)

	return module
}
