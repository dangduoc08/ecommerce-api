package users

import (
	"github.com/dangduoc08/gooh/core"
)

var UserModule = core.ModuleBuilder().
	Controllers(
		UserController{},
	).
	Providers(
		UserProvider{},
	).
	Exports(
		UserProvider{},
	).
	Build()
