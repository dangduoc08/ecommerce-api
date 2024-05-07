package auths

import (
	"github.com/dangduoc08/ecommerce-api/admins/users"
	"github.com/dangduoc08/ecommerce-api/mails"
	"github.com/dangduoc08/gogo/core"
)

var AuthModule = func() *core.Module {
	module := core.ModuleBuilder().
		Imports(
			mails.MailModule,
		).
		Controllers(
			AuthController{},
		).
		Providers(
			AuthProvider{},
			users.UserProvider{},
		).
		Build()

	module.
		Prefix("auths")

	return module
}()
