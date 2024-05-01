package users

import (
	"github.com/dangduoc08/gogo/core"
)

var UserModule = func() *core.Module {
	module := core.ModuleBuilder().
		Controllers(
			UserController{},
		).
		Providers(
			UserProvider{},
		).
		Build()

	module.
		Prefix("users")

	return module
}()
